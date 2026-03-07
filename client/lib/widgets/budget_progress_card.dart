import 'package:flutter/material.dart';
import '../models/models.dart';
import '../utils/utils.dart';

class BudgetProgressCard extends StatelessWidget {
  final ProjectBudget budget;
  final VoidCallback? onTap;
  final VoidCallback? onEdit;

  const BudgetProgressCard({
    super.key,
    required this.budget,
    this.onTap,
    this.onEdit,
  });

  double get percentage =>
      budget.budget > 0 ? (budget.used / budget.budget).clamp(0.0, 1.0) : 0;

  Color _getColorByPercentage() {
    if (percentage >= 0.9) return AppTheme.error;
    if (percentage >= 0.7) return AppTheme.warning;
    return AppTheme.success;
  }

  @override
  Widget build(BuildContext context) {
    final color = _getColorByPercentage();

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
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '预算 #${budget.id}',
                    style: Theme.of(context).textTheme.titleSmall,
                  ),
                  if (onEdit != null)
                    IconButton(
                      icon: const Icon(Icons.edit_outlined),
                      iconSize: 18,
                      onPressed: onEdit,
                      tooltip: '编辑预算',
                    ),
                ],
              ),
              const SizedBox(height: 12),
              LinearBatteryGauge(
                percentage: percentage,
                height: 12,
                fillColor: color,
                showLabel: false,
              ),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '已使用: ${CurrencyUtils.format(budget.used)}',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          color: Colors.grey.shade600,
                        ),
                  ),
                  Text(
                    '预算: ${CurrencyUtils.format(budget.budget)}',
                    style: Theme.of(context).textTheme.bodySmall?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                ],
              ),
              const SizedBox(height: 4),
              Text(
                '剩余: ${CurrencyUtils.format(budget.budget - budget.used)}',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: color,
                      fontWeight: FontWeight.bold,
                    ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
