// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'account.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$AccountImpl _$$AccountImplFromJson(Map<String, dynamic> json) =>
    _$AccountImpl(
      id: (json['id'] as num).toInt(),
      userId: (json['userId'] as num).toInt(),
      type: json['type'] as String,
      balance: (json['balance'] as num).toDouble(),
      version: (json['version'] as num).toInt(),
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$$AccountImplToJson(_$AccountImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'userId': instance.userId,
      'type': instance.type,
      'balance': instance.balance,
      'version': instance.version,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
    };

_$CreateAccountRequestImpl _$$CreateAccountRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateAccountRequestImpl(
      type: json['type'] as String,
      balance: (json['balance'] as num?)?.toDouble(),
    );

Map<String, dynamic> _$$CreateAccountRequestImplToJson(
        _$CreateAccountRequestImpl instance) =>
    <String, dynamic>{
      'type': instance.type,
      'balance': instance.balance,
    };

_$AdjustBalanceRequestImpl _$$AdjustBalanceRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$AdjustBalanceRequestImpl(
      amount: (json['amount'] as num).toDouble(),
    );

Map<String, dynamic> _$$AdjustBalanceRequestImplToJson(
        _$AdjustBalanceRequestImpl instance) =>
    <String, dynamic>{
      'amount': instance.amount,
    };
