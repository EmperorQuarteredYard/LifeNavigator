import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import '../utils/utils.dart';

class BatteryGauge extends StatelessWidget {
  final double percentage;
  final double? width;
  final double? height;
  final Color? fillColor;
  final Color? backgroundColor;
  final bool showLabel;
  final String? label;

  const BatteryGauge({
    super.key,
    required this.percentage,
    this.width,
    this.height,
    this.fillColor,
    this.backgroundColor,
    this.showLabel = true,
    this.label,
  });

  Color _getColorByPercentage() {
    if (percentage >= 0.7) return AppTheme.success;
    if (percentage >= 0.3) return AppTheme.warning;
    return AppTheme.error;
  }

  @override
  Widget build(BuildContext context) {
    final color = fillColor ?? _getColorByPercentage();
    final bgColor = backgroundColor ?? Colors.grey.shade300;

    return SizedBox(
      width: width ?? 150,
      height: height ?? 80,
      child: Stack(
        alignment: Alignment.center,
        children: [
          AspectRatio(
            aspectRatio: 2,
            child: PieChart(
              PieChartData(
                startAngleOffset: -90,
                sectionsSpace: 0,
                centerSpaceRadius: 30,
                sections: [
                  PieChartSectionData(
                    value: percentage.clamp(0.0, 1.0),
                    color: color,
                    radius: 15,
                    title: '',
                    showTitle: false,
                  ),
                  PieChartSectionData(
                    value: (1 - percentage).clamp(0.0, 1.0),
                    color: bgColor,
                    radius: 15,
                    title: '',
                    showTitle: false,
                  ),
                ],
              ),
            ),
          ),
          if (showLabel)
            Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  label ?? '${(percentage * 100).toStringAsFixed(0)}%',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.bold,
                        color: color,
                      ),
                ),
              ],
            ),
        ],
      ),
    );
  }
}

class LinearBatteryGauge extends StatelessWidget {
  final double percentage;
  final double? width;
  final double height;
  final Color? fillColor;
  final Color? backgroundColor;
  final bool showLabel;

  const LinearBatteryGauge({
    super.key,
    required this.percentage,
    this.width,
    this.height = 24,
    this.fillColor,
    this.backgroundColor,
    this.showLabel = true,
  });

  Color _getColorByPercentage() {
    if (percentage >= 0.7) return AppTheme.success;
    if (percentage >= 0.3) return AppTheme.warning;
    return AppTheme.error;
  }

  @override
  Widget build(BuildContext context) {
    final color = fillColor ?? _getColorByPercentage();
    final bgColor = backgroundColor ?? Colors.grey.shade300;

    return SizedBox(
      width: width,
      height: height,
      child: Stack(
        alignment: Alignment.centerLeft,
        children: [
          Container(
            height: height,
            decoration: BoxDecoration(
              color: bgColor,
              borderRadius: BorderRadius.circular(height / 2),
            ),
          ),
          FractionallySizedBox(
            widthFactor: percentage.clamp(0.0, 1.0),
            child: Container(
              height: height,
              decoration: BoxDecoration(
                color: color,
                borderRadius: BorderRadius.circular(height / 2),
              ),
            ),
          ),
          if (showLabel)
            Center(
              child: Text(
                '${(percentage * 100).toStringAsFixed(0)}%',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(
                      color: percentage > 0.5 ? Colors.white : Colors.black87,
                      fontWeight: FontWeight.bold,
                    ),
              ),
            ),
        ],
      ),
    );
  }
}
