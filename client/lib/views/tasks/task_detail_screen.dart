import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../models/models.dart';
import '../../providers/providers.dart';
import '../../widgets/widgets.dart';
import '../../utils/utils.dart';

class TaskDetailScreen extends ConsumerStatefulWidget {
  final int taskId;

  const TaskDetailScreen({super.key, required this.taskId});

  @override
  ConsumerState<TaskDetailScreen> createState() => _TaskDetailScreenState();
}

class _TaskDetailScreenState extends ConsumerState<TaskDetailScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(taskStateProvider.notifier).selectTask(widget.taskId);
    });
  }

  @override
  Widget build(BuildContext context) {
    final taskState = ref.watch(taskStateProvider);
    final task = taskState.selectedTask;

    if (taskState.isLoading && task == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('加载中...')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    if (task == null) {
      return Scaffold(
        appBar: AppBar(title: const Text('任务不存在')),
        body: const Center(child: Text('任务不存在或已被删除')),
      );
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(task.name),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit_outlined),
            onPressed: () => context.push('/tasks/${task.id}/edit'),
          ),
          IconButton(
            icon: const Icon(Icons.delete_outline),
            onPressed: () => _showDeleteDialog(task),
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildStatusCard(task),
            const SizedBox(height: 16),
            if (task.description != null && task.description!.isNotEmpty) ...[
              Text(
                '描述',
                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
              ),
              const SizedBox(height: 8),
              Text(task.description!),
              const SizedBox(height: 16),
            ],
            _buildInfoCard(task),
            const SizedBox(height: 16),
            _buildPaymentsSection(task),
          ],
        ),
      ),
      floatingActionButton: task.status != TaskStatus.completed.value
          ? FloatingActionButton.extended(
              onPressed: () => _handleComplete(task),
              icon: const Icon(Icons.check),
              label: const Text('完成任务'),
            )
          : null,
    );
  }

  Widget _buildStatusCard(Task task) {
    Color statusColor;
    String statusText;

    switch (TaskStatus.fromValue(task.status)) {
      case TaskStatus.pending:
        statusColor = AppTheme.pendingColor;
        statusText = '待处理';
        break;
      case TaskStatus.inProgress:
        statusColor = AppTheme.inProgressColor;
        statusText = '进行中';
        break;
      case TaskStatus.completed:
        statusColor = AppTheme.completedColor;
        statusText = '已完成';
        break;
      case TaskStatus.cancelled:
        statusColor = AppTheme.cancelledColor;
        statusText = '已取消';
        break;
    }

    return Card(
      color: statusColor.withOpacity(0.1),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(
          children: [
            Icon(
              _getStatusIcon(TaskStatus.fromValue(task.status)),
              color: statusColor,
              size: 32,
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    statusText,
                    style: Theme.of(context).textTheme.titleMedium?.copyWith(
                          color: statusColor,
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                  if (task.deadline != null)
                    Text(
                      DateTimeUtils.formatDeadline(task.deadline),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: DateTimeUtils.isOverdue(task.deadline)
                                ? AppTheme.error
                                : null,
                          ),
                    ),
                ],
              ),
            ),
            if (task.completedAt != null)
              Text(
                '完成于 ${DateTimeUtils.formatDate(task.completedAt)}',
                style: Theme.of(context).textTheme.bodySmall,
              ),
          ],
        ),
      ),
    );
  }

  IconData _getStatusIcon(TaskStatus status) {
    switch (status) {
      case TaskStatus.pending:
        return Icons.pending_outlined;
      case TaskStatus.inProgress:
        return Icons.autorenew;
      case TaskStatus.completed:
        return Icons.check_circle_outline;
      case TaskStatus.cancelled:
        return Icons.cancel_outlined;
    }
  }

  Widget _buildInfoCard(Task task) {
    return Card(
      child: Column(
        children: [
          ListTile(
            leading: const Icon(Icons.folder_outlined),
            title: const Text('所属项目'),
            subtitle: FutureBuilder<String>(
              future: _getProjectName(task.projectId),
              builder: (context, snapshot) {
                return Text(snapshot.data ?? '加载中...');
              },
            ),
            onTap: () => context.push('/projects/${task.projectId}'),
          ),
          ListTile(
            leading: const Icon(Icons.category_outlined),
            title: const Text('类型'),
            subtitle: Text(_getTypeText(TaskType.fromValue(task.type))),
          ),
          if (task.category != null && task.category!.isNotEmpty)
            ListTile(
              leading: const Icon(Icons.label_outlined),
              title: const Text('分类'),
              subtitle: Text(task.category!),
            ),
          ListTile(
            leading: const Icon(Icons.calendar_today_outlined),
            title: const Text('创建时间'),
            subtitle: Text(DateTimeUtils.formatDateTime(task.createdAt)),
          ),
          ListTile(
            leading: const Icon(Icons.update_outlined),
            title: const Text('更新时间'),
            subtitle: Text(DateTimeUtils.formatDateTime(task.updatedAt)),
          ),
        ],
      ),
    );
  }

  Future<String> _getProjectName(int projectId) async {
    final project = ref.read(projectByIdProvider(projectId));
    if (project != null) return project.name;

    try {
      final loadedProject =
          await ref.read(projectServiceProvider).getProject(projectId);
      return loadedProject.name;
    } catch (e) {
      return '未知项目';
    }
  }

  Widget _buildPaymentsSection(Task task) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              '付款记录',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                  ),
            ),
            TextButton.icon(
              onPressed: () => _showAddPaymentDialog(task),
              icon: const Icon(Icons.add),
              label: const Text('添加'),
            ),
          ],
        ),
        const SizedBox(height: 8),
        if (task.payments == null || task.payments!.isEmpty)
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  Icon(Icons.info_outline, color: Colors.grey.shade600),
                  const SizedBox(width: 12),
                  const Text('暂无付款记录'),
                ],
              ),
            ),
          )
        else
          ...task.payments!.map(
            (payment) => Card(
              child: ListTile(
                leading: const Icon(Icons.payment_outlined),
                title: Text(CurrencyUtils.format(payment.amount)),
                subtitle: Text('预算 #${payment.budgetId}'),
                trailing: IconButton(
                  icon: const Icon(Icons.delete_outline),
                  onPressed: () => _deletePayment(task, payment),
                ),
              ),
            ),
          ),
      ],
    );
  }

  Future<void> _handleComplete(Task task) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('完成任务'),
        content: Text('确定要将任务 "${task.name}" 标记为已完成吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            child: const Text('确定'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      await ref
          .read(taskStateProvider.notifier)
          .finishTask(task.id, DateTime.now());
    }
  }

  Future<void> _showDeleteDialog(Task task) async {
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

    if (confirmed == true && mounted) {
      await ref.read(taskStateProvider.notifier).deleteTask(task.id);
      if (mounted) Navigator.pop(context);
    }
  }

  Future<void> _showAddPaymentDialog(Task task) async {
    final amountController = TextEditingController();
    int? selectedBudgetId;

    final project = await ref
        .read(projectServiceProvider)
        .getProject(task.projectId);

    if (project.budgets == null || project.budgets!.isEmpty) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('该项目暂无预算，请先添加预算')),
        );
      }
      return;
    }

    await showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('添加付款'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              DropdownButtonFormField<int>(
                value: selectedBudgetId,
                decoration: const InputDecoration(
                  labelText: '选择预算',
                ),
                items: project.budgets!.map((budget) {
                  return DropdownMenuItem(
                    value: budget.id,
                    child: Text(
                      '预算 #${budget.id} - ${CurrencyUtils.format(budget.used)}/${CurrencyUtils.format(budget.budget)}',
                    ),
                  );
                }).toList(),
                onChanged: (value) {
                  setState(() {
                    selectedBudgetId = value;
                  });
                },
              ),
              const SizedBox(height: 16),
              TextField(
                controller: amountController,
                decoration: const InputDecoration(
                  labelText: '付款金额',
                  prefixText: '¥',
                ),
                keyboardType: TextInputType.number,
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
                final amount = double.tryParse(amountController.text);
                if (amount == null || amount <= 0) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('请输入有效的金额')),
                  );
                  return;
                }
                if (selectedBudgetId == null) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('请选择预算')),
                  );
                  return;
                }

                Navigator.pop(context);
                await ref.read(taskStateProvider.notifier).createPayment(
                      task.id,
                      CreateTaskPaymentRequest(
                        budgetId: selectedBudgetId!,
                        amount: amount,
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

  Future<void> _deletePayment(Task task, TaskPayment payment) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除付款记录'),
        content: const Text('确定要删除此付款记录吗？'),
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
      await ref.read(taskServiceProvider).deletePayment(task.id, payment.id);
      ref.read(taskStateProvider.notifier).selectTask(task.id);
    }
  }

  String _getTypeText(TaskType type) {
    switch (type) {
      case TaskType.normal:
        return '普通任务';
      case TaskType.milestone:
        return '里程碑';
      case TaskType.routine:
        return '日常任务';
    }
  }
}
