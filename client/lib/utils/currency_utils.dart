import 'package:intl/intl.dart';

class CurrencyUtils {
  static final NumberFormat _currencyFormat = NumberFormat.currency(
    locale: 'zh_CN',
    symbol: '¥',
    decimalDigits: 2,
  );

  static final NumberFormat _compactFormat = NumberFormat.compactCurrency(
    locale: 'zh_CN',
    symbol: '¥',
    decimalDigits: 1,
  );

  static String format(double? amount) {
    if (amount == null) return '¥0.00';
    return _currencyFormat.format(amount);
  }

  static String formatCompact(double? amount) {
    if (amount == null) return '¥0';
    return _compactFormat.format(amount);
  }

  static String formatWithSign(double? amount) {
    if (amount == null) return '¥0.00';
    if (amount > 0) {
      return '+${_currencyFormat.format(amount)}';
    }
    return _currencyFormat.format(amount);
  }

  static String formatPercentage(double? value, {int decimalDigits = 1}) {
    if (value == null) return '0%';
    return '${(value * 100).toStringAsFixed(decimalDigits)}%';
  }
}
