import 'package:freezed_annotation/freezed_annotation.dart';

part 'account.freezed.dart';
part 'account.g.dart';

@freezed
class Account with _$Account {
  const factory Account({
    required int id,
    required int userId,
    required String type,
    required double balance,
    required int version,
    required DateTime createdAt,
    required DateTime updatedAt,
  }) = _Account;

  factory Account.fromJson(Map<String, dynamic> json) =>
      _$AccountFromJson(json);
}

@freezed
class CreateAccountRequest with _$CreateAccountRequest {
  const factory CreateAccountRequest({
    required String type,
    double? balance,
  }) = _CreateAccountRequest;

  factory CreateAccountRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateAccountRequestFromJson(json);
}

@freezed
class AdjustBalanceRequest with _$AdjustBalanceRequest {
  const factory AdjustBalanceRequest({
    required double amount,
  }) = _AdjustBalanceRequest;

  factory AdjustBalanceRequest.fromJson(Map<String, dynamic> json) =>
      _$AdjustBalanceRequestFromJson(json);
}
