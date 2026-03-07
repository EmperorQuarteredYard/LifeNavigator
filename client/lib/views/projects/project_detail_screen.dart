import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/models.dart';
import '../../providers/providers.dart';
import '../../widgets/widgets.dart';
import '../../utils/utils.dart';

class ProjectDetailScreen extends ConsumerStatefulWidget {
  final int projectId;

  const ProjectDetailScreen({super.key, required this.projectId});

  @override
  ConsumerState<ProjectDetailScreen> createState() => _ProjectDetailScreenState();
}

class _ProjectDetailScreenState extends ConsumerState<ProjectDetailScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(projectStateProvider.notifier).selectProject(widget.projectId);
      ref.read(taskStateProvider.notifier).loadTasks(
            projectId: widget.projectId,
            refresh: true,
          );
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final projectState = ref.watch(projectStateProvider);
    final taskState = ref.watch(taskStateProvider);
    final project = projectState.selectedProject;

    if (projectState.isLoading && project == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('加载中...')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    if (project == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('项目不存在')),
        body: const Center(child: Text('项目不存在或已被删除')),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(project.name),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit_outlined),
            onPressed: () => context.push('/projects/${project.id}/edit'),
          ),
          IconButton(
            icon: const Icon(Icons.delete_outline),
            onPressed: () => _showDeleteDialog(project),
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '概览'),
            Tab(text: '任务'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildOverviewTab(project),
          _buildTasksTab(taskState),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/tasks/create?project_id=${project.id}'),
        child: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildOverviewTab(Project project) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (project.description != null && project.description!.isNotEmpty) ...[
            Text(
              '描述',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            const SizedBox(height: 8),
            Text(project.description!),
            const SizedBox(height: 24),
          ],
          Text(
            '预算',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 12),
          if (project.budgets == null || project.budgets!.isEmpty)
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  children: [
                    const Icon(Icons.info_outline, color: Colors.grey),
                    const SizedBox(width: 12),
                    const Text('暂无预算'),
                    const Spacer(),
                    TextButton(
                      onPressed: () => _showAddBudgetDialog(project),
                      child: const Text('添加预算'),
                    ),
                  ],
                ),
              ),
            )
          else
            ...project.budgets!.map(
              (budget) => Padding(
                padding: const EdgeInsets.only(bottom: 12),
                child: BudgetProgressCard(
                  budget: budget,
                  onEdit: () => _showEditBudgetDialog(project, budget),
                ),
              ),
            ),
          const SizedBox(height: 16),
          ElevatedButton.icon(
            onPressed: () => _showAddBudgetDialog(project),
            icon: const Icon(Icons.add),
            label: const Text('添加预算'),
          ),
          const SizedBox(height: 24),
          Text(
            '项目信息',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
                  fontWeight: FontWeight.bold,
                ),
          ),
          const SizedBox(height: 12),
          Card(
            child: Column(
              children: [
                ListTile(
                  leading: const Icon(Icons.calendar_today_outlined),
                  title: const Text('创建时间'),
                  subtitle: Text(DateTimeUtils.formatDateTime(project.createdAt)),
                ),
                ListTile(
                  leading: const Icon(Icons.update_outlined),
                  title: const Text('更新时间'),
                  subtitle: Text(DateTimeUtils.formatDateTime(project.updatedAt)),
                ),
                ListTile(
                  leading: const Icon(Icons.refresh_outlined),
                  title: const Text('刷新间隔'),
                  subtitle: Text(_getRefreshIntervalText(project.refreshInterval)),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTasksTab(TaskState taskState) {
    final tasks = taskState.tasks
        .where((t) => t.projectId == widget.projectId)
        .toList();

    if (taskState.isLoading && tasks.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (tasks.isEmpty) {
      return const EmptyState(
        title: '暂无任务',
        subtitle: '点击右下角按钮创建新任务',
        icon: Icons.task_outlined,
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: tasks.length,
      itemBuilder: (context, index) {
        final task = tasks[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: TaskCard(
            task: task,
            onTap: () => context.push('/tasks/${task.id}'),
            onComplete: () async {
              await ref
                  .read(taskStateProvider.notifier)
                  .finishTask(task.id, DateTime.now());
            },
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
    );
  }

  String _getRefreshIntervalText(int interval) {
    switch (interval) {
      case 0:
        return '从不';
      case 1:
        return '每年';
      case 2:
        return '每月';
      case 3:
        return '每周';
      case 4:
        return '每天';
      case 5:
        return '每小时';
      default:
        return '未知';
    }
  }

  Future<void> _showDeleteDialog(Project project) async {
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
      if (mounted) Navigator.pop(context);
    }
  }

  Future<void> _showAddBudgetDialog(Project project) async {
    final budgetController = TextEditingController();
    int? selectedAccountId;

    await showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('添加预算'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: budgetController,
                decoration: const InputDecoration(
                  labelText: '预算金额',
                  prefixText: '¥',
                ),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<int>(
                value: selectedAccountId,
                decoration: const InputDecoration(
                  labelText: '关联账户（可选）',
                ),
                items: [
                  const DropdownMenuItem(value: null, child: Text('不关联')),
                  ...ref.read(accountStateProvider).accounts.map(
                        (account) => DropdownMenuItem(
                          value: account.id,
                          child: Text('${account.type} - ${CurrencyUtils.format(account.balance)}'),
                        ),
                      ),
                ],
                onChanged: (value) {
                  setState(() {
                    selectedAccountId = value;
                  });
                },
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('取消'),
            ),
            TextButton(
              onPressed: () async {
                final budget = double.tryParse(budgetController.text);
                if (budget == null || budget <= 0) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('请输入有效的预算金额')),
                  );
                  return;
                }

                Navigator.pop(context);
                await ref.read(projectStateProvider.notifier).createBudget(
                      project.id,
                      ProjectBudgetRequest(
                        budget: budget,
                        accountId: selectedAccountId,
                      ),
                    );
              },
              child: const Text('添加'),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _showEditBudgetDialog(Project project, ProjectBudget budget) async {
    final budgetController = TextEditingController(text: budget.budget.toString());
    int? selectedAccountId = budget.accountId;

    await showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('编辑预算'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: budgetController,
                decoration: const InputDecoration(
                  labelText: '预算金额',
                  prefixText: '¥',
                ),
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<int>(
                value: selectedAccountId,
                decoration: const InputDecoration(
                  labelText: '关联账户（可选）',
                ),
                items: [
                  const DropdownMenuItem(value: null, child: Text('不关联')),
                  ...ref.read(accountStateProvider).accounts.map(
                        (account) => DropdownMenuItem(
                          value: account.id,
                          child: Text('${account.type} - ${CurrencyUtils.format(account.balance)}'),
                        ),
                      ),
                ],
                onChanged: (value) {
                  setState(() {
                    selectedAccountId = value;
                  });
                },
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('取消'),
            ),
            TextButton(
              onPressed: () async {
                final newBudget = double.tryParse(budgetController.text);
                if (newBudget == null || newBudget <= 0) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('请输入有效的预算金额')),
                  );
                  return;
                }

                Navigator.pop(context);
                await ref.read(projectServiceProvider).updateBudget(
                      project.id,
                      budget.id,
                      ProjectBudgetRequest(
                        budget: newBudget,
                        accountId: selectedAccountId,
                      ),
                    );
                ref.read(projectStateProvider.notifier).selectProject(project.id);
              },
              child: const Text('保存'),
            ),
          ],
        ),
      ),
    );
  }
}
