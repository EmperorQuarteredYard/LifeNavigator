import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/models.dart';
import '../../providers/providers.dart';

class TaskFormScreen extends ConsumerStatefulWidget {
  final int? taskId;
  final int? projectId;

  const TaskFormScreen({super.key, this.taskId, this.projectId});

  @override
  ConsumerState<TaskFormScreen> createState() => _TaskFormScreenState();
}

class _TaskFormScreenState extends ConsumerState<TaskFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _categoryController = TextEditingController();
  int? _selectedProjectId;
  TaskStatus _selectedStatus = TaskStatus.pending;
  TaskType _selectedType = TaskType.normal;
  DateTime? _deadline;
  bool _autoCalculated = false;

  bool get isEditing => widget.taskId != null;

  @override
  void initState() {
    super.initState();
    _selectedProjectId = widget.projectId;
    if (isEditing) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final task = ref.read(taskByIdProvider(widget.taskId!));
        if (task != null) {
          _nameController.text = task.name;
          _descriptionController.text = task.description ?? '';
          _categoryController.text = task.category ?? '';
          _selectedProjectId = task.projectId;
          _selectedStatus = TaskStatus.fromValue(task.status);
          _selectedType = TaskType.fromValue(task.type);
          _deadline = task.deadline;
          _autoCalculated = task.autoCalculated;
        }
      });
    }
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(projectStateProvider.notifier).loadProjects(refresh: true);
    });
  }

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    _categoryController.dispose();
    super.dispose();
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    final request = CreateTaskRequest(
      projectId: _selectedProjectId!,
      name: _nameController.text.trim(),
      description: _descriptionController.text.trim().isNotEmpty
          ? _descriptionController.text.trim()
          : null,
      autoCalculated: _autoCalculated,
      type: _selectedType.value,
      status: _selectedStatus.value,
      category: _categoryController.text.trim().isNotEmpty
          ? _categoryController.text.trim()
          : null,
      deadline: _deadline?.toIso8601String(),
    );

    bool success;
    if (isEditing) {
      success = await ref.read(taskStateProvider.notifier).updateTask(
            widget.taskId!,
            UpdateTaskRequest(
              name: request.name,
              projectId: request.projectId,
              description: request.description,
              autoCalculated: request.autoCalculated,
              type: request.type,
              status: request.status,
              category: request.category,
              deadline: request.deadline,
            ),
          );
    } else {
      success = await ref.read(taskStateProvider.notifier).createTask(request);
    }

    if (success && mounted) {
      Navigator.pop(context);
    }
  }

  Future<void> _selectDeadline() async {
    final date = await showDatePicker(
      context: context,
      initialDate: _deadline ?? DateTime.now(),
      firstDate: DateTime.now().subtract(const Duration(days: 365)),
      lastDate: DateTime.now().add(const Duration(days: 365 * 5)),
    );

    if (date != null && mounted) {
      final time = await showTimePicker(
        context: context,
        initialTime: TimeOfDay.fromDateTime(_deadline ?? DateTime.now()),
      );

      if (time != null) {
        setState(() {
          _deadline = DateTime(
            date.year,
            date.month,
            date.day,
            time.hour,
            time.minute,
          );
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final projectState = ref.watch(projectStateProvider);
    final taskState = ref.watch(taskStateProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(isEditing ? '编辑任务' : '新建任务'),
        actions: [
          TextButton(
            onPressed: (taskState.isLoading || _selectedProjectId == null)
                ? null
                : _handleSubmit,
            child: const Text('保存'),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            DropdownButtonFormField<int>(
              value: _selectedProjectId,
              decoration: const InputDecoration(
                labelText: '所属项目 *',
                prefixIcon: Icon(Icons.folder_outlined),
              ),
              items: projectState.projects.map((project) {
                return DropdownMenuItem(
                  value: project.id,
                  child: Text(project.name),
                );
              }).toList(),
              validator: (value) {
                if (value == null) {
                  return '请选择项目';
                }
                return null;
              },
              onChanged: (value) {
                setState(() {
                  _selectedProjectId = value;
                });
              },
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: '任务名称 *',
                prefixIcon: Icon(Icons.task_outlined),
              ),
              textInputAction: TextInputAction.next,
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return '请输入任务名称';
                }
                return null;
              },
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _descriptionController,
              decoration: const InputDecoration(
                labelText: '任务描述',
                prefixIcon: Icon(Icons.description_outlined),
              ),
              maxLines: 3,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<TaskStatus>(
                    value: _selectedStatus,
                    decoration: const InputDecoration(
                      labelText: '状态',
                      prefixIcon: Icon(Icons.flag_outlined),
                    ),
                    items: TaskStatus.values.map((status) {
                      return DropdownMenuItem(
                        value: status,
                        child: Text(_getStatusText(status)),
                      );
                    }).toList(),
                    onChanged: (value) {
                      setState(() {
                        _selectedStatus = value ?? TaskStatus.pending;
                      });
                    },
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: DropdownButtonFormField<TaskType>(
                    value: _selectedType,
                    decoration: const InputDecoration(
                      labelText: '类型',
                      prefixIcon: Icon(Icons.category_outlined),
                    ),
                    items: TaskType.values.map((type) {
                      return DropdownMenuItem(
                        value: type,
                        child: Text(_getTypeText(type)),
                      );
                    }).toList(),
                    onChanged: (value) {
                      setState(() {
                        _selectedType = value ?? TaskType.normal;
                      });
                    },
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _categoryController,
              decoration: const InputDecoration(
                labelText: '分类',
                prefixIcon: Icon(Icons.label_outlined),
              ),
            ),
            const SizedBox(height: 16),
            ListTile(
              leading: const Icon(Icons.calendar_today_outlined),
              title: const Text('截止时间'),
              subtitle: Text(
                _deadline != null
                    ? '${_deadline!.year}-${_deadline!.month.toString().padLeft(2, '0')}-${_deadline!.day.toString().padLeft(2, '0')} ${_deadline!.hour.toString().padLeft(2, '0')}:${_deadline!.minute.toString().padLeft(2, '0')}'
                    : '未设置',
              ),
              trailing: _deadline != null
                  ? IconButton(
                      icon: const Icon(Icons.clear),
                      onPressed: () {
                        setState(() {
                          _deadline = null;
                        });
                      },
                    )
                  : null,
              onTap: _selectDeadline,
            ),
            const SizedBox(height: 8),
            SwitchListTile(
              secondary: const Icon(Icons.calculate_outlined),
              title: const Text('自动计算'),
              subtitle: const Text('自动计算任务相关数据'),
              value: _autoCalculated,
              onChanged: (value) {
                setState(() {
                  _autoCalculated = value;
                });
              },
            ),
            if (taskState.error != null) ...[
              const SizedBox(height: 16),
              Text(
                taskState.error!,
                style: TextStyle(
                  color: Theme.of(context).colorScheme.error,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }

  String _getStatusText(TaskStatus status) {
    switch (status) {
      case TaskStatus.pending:
        return '待处理';
      case TaskStatus.inProgress:
        return '进行中';
      case TaskStatus.completed:
        return '已完成';
      case TaskStatus.cancelled:
        return '已取消';
    }
  }

  String _getTypeText(TaskType type) {
    switch (type) {
      case TaskType.normal:
        return '普通';
      case TaskType.milestone:
        return '里程碑';
      case TaskType.routine:
        return '日常';
    }
  }
}
