import 'package:flutter/material.dart';
import '../models/models.dart';
import '../utils/utils.dart';

class ProjectCard extends StatelessWidget {
  final Project project;
  final VoidCallback? onTap;
  final VoidCallback? onEdit;
  final VoidCallback? onDelete;

  const ProjectCard({
    super.key,
    required this.project,
    this.onTap,
    this.onEdit,
    this.onDelete,
  });

  double _getBudgetPercentage() {
    if (project.budgets == null || project.budgets!.isEmpty) {
      return 0;
    }
    final totalBudget = project.budgets!.fold<double>(
      0,
      (sum, budget) => sum + budget.budget,
    );
    final totalUsed = project.budgets!.fold<double>(
      0,
      (sum, budget) => sum + budget.used,
    );
    if (totalBudget == 0) return 0;
    return totalUsed / totalBudget;
  }

  @override
  Widget build(BuildContext context) {
    final budgetPercentage = _getBudgetPercentage();

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
                      project.name,
                      style: Theme.of(context).textTheme.titleMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                    ),
                  ),
                  if (onEdit != null)
                    IconButton(
                      icon: const Icon(Icons.edit_outlined),
                      iconSize: 20,
                      onPressed: onEdit,
                      tooltip: '编辑项目',
                    ),
                  if (onDelete != null)
                    IconButton(
                      icon: const Icon(Icons.delete_outline),
                      iconSize: 20,
                      color: AppTheme.error,
                      onPressed: onDelete,
                      tooltip: '删除项目',
                    ),
                ],
              ),
              if (project.description != null &&
                  project.description!.isNotEmpty) ...[
                const SizedBox(height: 8),
                Text(
                  project.description!,
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: Colors.grey.shade600,
                      ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
              if (project.budgets != null && project.budgets!.isNotEmpty) ...[
                const SizedBox(height: 12),
                LinearBatteryGauge(
                  percentage: budgetPercentage.clamp(0.0, 1.0),
                  height: 8,
                  showLabel: false,
                ),
                const SizedBox(height: 4),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '预算使用',
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            color: Colors.grey.shade600,
                          ),
                    ),
                    Text(
                      CurrencyUtils.format(
                        project.budgets!.fold<double>(
                          0,
                          (sum, budget) => sum + budget.used,
                        ),
                      ),
                      style: Theme.of(context).textTheme.bodySmall?.copyWith(
                            fontWeight: FontWeight.bold,
                          ),
                    ),
                  ],
                ),
              ],
              const SizedBox(height: 8),
              Row(
                children: [
                  Icon(
                    Icons.calendar_today_outlined,
                    size: 14,
                    color: Colors.grey.shade600,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    DateTimeUtils.formatDate(project.createdAt),
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Colors.grey.shade600,
                        ),
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
