import 'package:intl/intl.dart';

class DateTimeUtils {
  static final DateFormat _dateFormat = DateFormat('yyyy-MM-dd');
  static final DateFormat _timeFormat = DateFormat('HH:mm');
  static final DateFormat _dateTimeFormat = DateFormat('yyyy-MM-dd HH:mm');
  static final DateFormat _monthFormat = DateFormat('yyyy-MM');

  static String formatDate(DateTime? date) {
    if (date == null) return '';
    return _dateFormat.format(date);
  }

  static String formatTime(DateTime? date) {
    if (date == null) return '';
    return _timeFormat.format(date);
  }

  static String formatDateTime(DateTime? date) {
    if (date == null) return '';
    return _dateTimeFormat.format(date);
  }

  static String formatMonth(DateTime? date) {
    if (date == null) return '';
    return _monthFormat.format(date);
  }

  static String formatRelative(DateTime? date) {
    if (date == null) return '';
    final now = DateTime.now();
    final diff = date.difference(now);

    if (diff.isNegative) {
      final absDiff = diff.abs();
      if (absDiff.inDays == 0) {
        return '今天';
      } else if (absDiff.inDays == 1) {
        return '昨天';
      } else if (absDiff.inDays < 7) {
        return '${absDiff.inDays}天前';
      } else {
        return formatDate(date);
      }
    } else {
      if (diff.inDays == 0) {
        return '今天';
      } else if (diff.inDays == 1) {
        return '明天';
      } else if (diff.inDays < 7) {
        return '${diff.inDays}天后';
      } else {
        return formatDate(date);
      }
    }
  }

  static String formatDeadline(DateTime? deadline) {
    if (deadline == null) return '无截止日期';

    final now = DateTime.now();
    final diff = deadline.difference(now);

    if (diff.isNegative) {
      return '已过期 ${formatDate(deadline)}';
    } else if (diff.inDays == 0) {
      return '今天截止';
    } else if (diff.inDays == 1) {
      return '明天截止';
    } else if (diff.inDays < 7) {
      return '${diff.inDays}天后截止';
    } else {
      return formatDate(deadline);
    }
  }

  static bool isOverdue(DateTime? deadline) {
    if (deadline == null) return false;
    return deadline.isBefore(DateTime.now());
  }

  static bool isToday(DateTime? date) {
    if (date == null) return false;
    final now = DateTime.now();
    return date.year == now.year &&
        date.month == now.month &&
        date.day == now.day;
  }

  static bool isTomorrow(DateTime? date) {
    if (date == null) return false;
    final tomorrow = DateTime.now().add(const Duration(days: 1));
    return date.year == tomorrow.year &&
        date.month == tomorrow.month &&
        date.day == tomorrow.day;
  }
}
