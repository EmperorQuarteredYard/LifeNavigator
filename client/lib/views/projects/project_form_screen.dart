import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/models.dart';
import '../../providers/providers.dart';

class ProjectFormScreen extends ConsumerStatefulWidget {
  final int? projectId;

  const ProjectFormScreen({super.key, this.projectId});

  @override
  ConsumerState<ProjectFormScreen> createState() => _ProjectFormScreenState();
}

class _ProjectFormScreenState extends ConsumerState<ProjectFormScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _descriptionController = TextEditingController();
  int _refreshInterval = 0;

  bool get isEditing => widget.projectId != null;

  @override
  void initState() {
    super.initState();
    if (isEditing) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final project = ref.read(projectByIdProvider(widget.projectId!));
        if (project != null) {
          _nameController.text = project.name;
          _descriptionController.text = project.description ?? '';
          _refreshInterval = project.refreshInterval;
        }
      });
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  Future<void> _handleSubmit() async {
    if (!_formKey.currentState!.validate()) return;

    final request = CreateProjectRequest(
      name: _nameController.text.trim(),
      description: _descriptionController.text.trim().isNotEmpty
          ? _descriptionController.text.trim()
          : null,
      refreshInterval: _refreshInterval,
    );

    bool success;
    if (isEditing) {
      success = await ref.read(projectStateProvider.notifier).updateProject(
            widget.projectId!,
            UpdateProjectRequest(
              name: request.name!,
              description: request.description,
              refreshInterval: request.refreshInterval,
            ),
          );
    } else {
      success = await ref.read(projectStateProvider.notifier).createProject(request);
    }

    if (success && mounted) {
      Navigator.pop(context);
    }
  }

  @override
  Widget build(BuildContext context) {
    final projectState = ref.watch(projectStateProvider);

    return Scaffold(
      appBar: AppBar(
        title: Text(isEditing ? '编辑项目' : '新建项目'),
        actions: [
          TextButton(
            onPressed: projectState.isLoading ? null : _handleSubmit,
            child: const Text('保存'),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            TextFormField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: '项目名称 *',
                prefixIcon: Icon(Icons.folder_outlined),
              ),
              textInputAction: TextInputAction.next,
              validator: (value) {
                if (value == null || value.trim().isEmpty) {
                  return '请输入项目名称';
                }
                return null;
              },
            ),
            const SizedBox(height: 16),
            TextFormField(
              controller: _descriptionController,
              decoration: const InputDecoration(
                labelText: '项目描述',
                prefixIcon: Icon(Icons.description_outlined),
              ),
              maxLines: 3,
              textInputAction: TextInputAction.newline,
            ),
            const SizedBox(height: 16),
            DropdownButtonFormField<int>(
              value: _refreshInterval,
              decoration: const InputDecoration(
                labelText: '刷新间隔',
                prefixIcon: Icon(Icons.refresh_outlined),
              ),
              items: const [
                DropdownMenuItem(value: 0, child: Text('从不')),
                DropdownMenuItem(value: 1, child: Text('每年')),
                DropdownMenuItem(value: 2, child: Text('每月')),
                DropdownMenuItem(value: 3, child: Text('每周')),
                DropdownMenuItem(value: 4, child: Text('每天')),
                DropdownMenuItem(value: 5, child: Text('每小时')),
              ],
              onChanged: (value) {
                setState(() {
                  _refreshInterval = value ?? 0;
                });
              },
            ),
            if (projectState.error != null) ...[
              const SizedBox(height: 16),
              Text(
                projectState.error!,
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
}
