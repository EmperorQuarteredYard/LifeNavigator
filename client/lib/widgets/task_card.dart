import 'package:flutter/material.dart';
import '../models/models.dart';
import '../utils/utils.dart';

class TaskCard extends StatelessWidget {
  final Task task;
  final VoidCallback? onTap;
  final VoidCallback? onComplete;
  final VoidCallback? onDelete;
  final bool showProject;

  const TaskCard({
    super.key,
    required this.task,
    this.onTap,
    this.onComplete,
    this.onDelete,
    this.showProject = false,
  });

  Color _getStatusColor() {
    switch (TaskStatus.fromValue(task.status)) {
      case TaskStatus.pending:
        return AppTheme.pendingColor;
      case TaskStatus.inProgress:
        return AppTheme.inProgressColor;
      case TaskStatus.completed:
        return AppTheme.completedColor;
      case TaskStatus.cancelled:
        return AppTheme.cancelledColor;
    }
  }

  String _getStatusText() {
    switch (TaskStatus.fromValue(task.status)) {
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

  @override
  Widget build(BuildContext context) {
    final isOverdue = DateTimeUtils.isOverdue(task.deadline) &&
        task.status != TaskStatus.completed.value;

    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      task.name,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            decoration: task.status == TaskStatus.completed.value
                                ? TextDecoration.lineThrough
                                : null,
                          ),
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: _getStatusColor().withOpacity(0.1),
                      borderRadius: BorderRadius.circular(4),
                    ),
                    child: Text(
                      _getStatusText(),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: _getStatusColor(),
                          ),
                    ),
                  ),
                ],
              ),
              if (task.description != null && task.description!.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text(
                  task.description!,
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: Colors.grey.shade600,
                      ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
              const SizedBox(height: 12),
              Row(
                children: [
                  if (task.deadline != null) ...[
                    Icon(
                      Icons.schedule,
                      size: 16,
                      color: isOverdue ? AppTheme.error : Colors.grey,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      DateTimeUtils.formatDeadline(task.deadline),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: isOverdue ? AppTheme.error : Colors.grey,
                          ),
                    ),
                    const SizedBox(width: 16),
                  ],
                  if (task.category != null && task.category!.isNotEmpty) ...[
                    Icon(
                      Icons.label_outline,
                      size: 16,
                      color: Colors.grey.shade600,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      task.category!,
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: Colors.grey.shade600,
                          ),
                    ),
                  ],
                  const Spacer(),
                  if (onComplete != null && task.status != TaskStatus.completed.value)
                    IconButton(
                      icon: const Icon(Icons.check_circle_outline),
                      iconSize: 20,
                      color: AppTheme.success,
                      onPressed: onComplete,
                      tooltip: '完成任务',
                    ),
                  if (onDelete != null)
                    IconButton(
                      icon: const Icon(Icons.delete_outline),
                      iconSize: 20,
                      color: AppTheme.error,
                      onPressed: onDelete,
                      tooltip: '删除任务',
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
