import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/models.dart';
import '../../providers/providers.dart';
import '../../widgets/widgets.dart';

class TaskListScreen extends ConsumerStatefulWidget {
  final int? projectId;

  const TaskListScreen({super.key, this.projectId});

  @override
  ConsumerState<TaskListScreen> createState() => _TaskListScreenState();
}

class _TaskListScreenState extends ConsumerState<TaskListScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  TaskStatus _selectedStatus = TaskStatus.pending;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    _tabController.addListener(_onTabChanged);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(taskStateProvider.notifier).loadTasks(
            projectId: widget.projectId,
            refresh: true,
          );
    });
  }

  @override
  void dispose() {
    _tabController.removeListener(_onTabChanged);
    _tabController.dispose();
    super.dispose();
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      setState(() {
        _selectedStatus = TaskStatus.values[_tabController.index];
      });
      ref.read(taskStateProvider.notifier).loadTasks(
            projectId: widget.projectId,
            refresh: true,
            status: _selectedStatus.value,
          );
    }
  }

  @override
  Widget build(BuildContext context) {
    final taskState = ref.watch(taskStateProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('任务列表'),
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '待处理'),
            Tab(text: '进行中'),
            Tab(text: '已完成'),
            Tab(text: '已取消'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildTaskList(taskState, TaskStatus.pending),
          _buildTaskList(taskState, TaskStatus.inProgress),
          _buildTaskList(taskState, TaskStatus.completed),
          _buildTaskList(taskState, TaskStatus.cancelled),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          final uri = widget.projectId != null
              ? '/tasks/create?project_id=${widget.projectId}'
              : '/tasks/create';
          context.push(uri);
        },
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildTaskList(TaskState taskState, TaskStatus status) {
    final tasks = taskState.tasks.where((t) => t.status == status.value).toList();

    if (taskState.isLoading && tasks.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (tasks.isEmpty) {
      return EmptyState(
        title: '暂无${_getStatusText(status)}任务',
        icon: Icons.task_outlined,
      );
    }

    return RefreshIndicator(
      onRefresh: () async {
        await ref.read(taskStateProvider.notifier).loadTasks(
              projectId: widget.projectId,
              refresh: true,
              status: status.value,
            );
      },
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: tasks.length,
        itemBuilder: (context, index) {
          final task = tasks[index];
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: TaskCard(
              task: task,
              showProject: widget.projectId == null,
              onTap: () => context.push('/tasks/${task.id}'),
              onComplete: status != TaskStatus.completed
                  ? () async {
                      await ref
                          .read(taskStateProvider.notifier)
                          .finishTask(task.id, DateTime.now());
                    }
                  : null,
              onDelete: () async {
                final confirmed = await showDialog<bool>(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('删除任务'),
                    content: Text('确定要删除任务 "${task.name}" 吗？'),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.pop(context, false),
                        child: const Text('取消'),
                      ),
                      TextButton(
                        onPressed: () => Navigator.pop(context, true),
                        style: TextButton.styleFrom(
                          foregroundColor: Theme.of(context).colorScheme.error,
                        ),
                        child: const Text('删除'),
                      ),
                    ],
                  ),
                );

                if (confirmed == true) {
                  await ref
                      .read(taskStateProvider.notifier)
                      .deleteTask(task.id);
                }
              },
            ),
          );
        },
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
}
