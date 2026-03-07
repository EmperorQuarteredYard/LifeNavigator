import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/providers.dart';
import '../../widgets/widgets.dart';

class ProjectListScreen extends ConsumerStatefulWidget {
  const ProjectListScreen({super.key});

  @override
  ConsumerState<ProjectListScreen> createState() => _ProjectListScreenState();
}

class _ProjectListScreenState extends ConsumerState<ProjectListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(projectStateProvider.notifier).loadProjects(refresh: true);
    });
  }

  @override
  Widget build(BuildContext context) {
    final projectState = ref.watch(projectStateProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('项目列表'),
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await ref.read(projectStateProvider.notifier).loadProjects(refresh: true);
        },
        child: LoadingOverlay(
          isLoading: projectState.isLoading && projectState.projects.isEmpty,
          child: projectState.projects.isEmpty
              ? const EmptyState(
                  title: '暂无项目',
                  subtitle: '点击右下角按钮创建新项目',
                  icon: Icons.folder_outlined,
                )
              : ListView.builder(
                  padding: const EdgeInsets.all(16),
                  itemCount: projectState.projects.length,
                  itemBuilder: (context, index) {
                    final project = projectState.projects[index];
                    return Padding(
                      padding: const EdgeInsets.only(bottom: 12),
                      child: ProjectCard(
                        project: project,
                        onTap: () => context.push('/projects/${project.id}'),
                        onEdit: () => context.push('/projects/${project.id}/edit'),
                        onDelete: () => _showDeleteDialog(project),
                      ),
                    );
                  },
                ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/projects/create'),
        child: const Icon(Icons.add),
      ),
    );
  }

  Future<void> _showDeleteDialog(project) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除项目'),
        content: Text('确定要删除项目 "${project.name}" 吗？\n此操作将同时删除项目下的所有任务。'),
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

    if (confirmed == true && mounted) {
      await ref.read(projectStateProvider.notifier).deleteProject(project.id);
    }
  }
}
