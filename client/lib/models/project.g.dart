// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'project.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$ProjectBudgetImpl _$$ProjectBudgetImplFromJson(Map<String, dynamic> json) =>
    _$ProjectBudgetImpl(
      id: (json['id'] as num).toInt(),
      projectId: (json['projectId'] as num).toInt(),
      accountId: (json['accountId'] as num?)?.toInt(),
      budget: (json['budget'] as num).toDouble(),
      used: (json['used'] as num).toDouble(),
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );

Map<String, dynamic> _$$ProjectBudgetImplToJson(_$ProjectBudgetImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'projectId': instance.projectId,
      'accountId': instance.accountId,
      'budget': instance.budget,
      'used': instance.used,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
    };

_$ProjectImpl _$$ProjectImplFromJson(Map<String, dynamic> json) =>
    _$ProjectImpl(
      id: (json['id'] as num).toInt(),
      userId: (json['userId'] as num).toInt(),
      name: json['name'] as String,
      description: json['description'] as String?,
      refreshInterval: (json['refreshInterval'] as num).toInt(),
      lastRefresh: DateTime.parse(json['lastRefresh'] as String),
      maxTaskId: (json['maxTaskId'] as num).toInt(),
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
      budgets: (json['budgets'] as List<dynamic>?)
          ?.map((e) => ProjectBudget.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$$ProjectImplToJson(_$ProjectImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'userId': instance.userId,
      'name': instance.name,
      'description': instance.description,
      'refreshInterval': instance.refreshInterval,
      'lastRefresh': instance.lastRefresh.toIso8601String(),
      'maxTaskId': instance.maxTaskId,
      'createdAt': instance.createdAt.toIso8601String(),
      'updatedAt': instance.updatedAt.toIso8601String(),
      'budgets': instance.budgets,
    };

_$CreateProjectRequestImpl _$$CreateProjectRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateProjectRequestImpl(
      name: json['name'] as String,
      description: json['description'] as String?,
      refreshInterval: (json['refreshInterval'] as num?)?.toInt(),
    );

Map<String, dynamic> _$$CreateProjectRequestImplToJson(
        _$CreateProjectRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'description': instance.description,
      'refreshInterval': instance.refreshInterval,
    };

_$UpdateProjectRequestImpl _$$UpdateProjectRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateProjectRequestImpl(
      name: json['name'] as String,
      description: json['description'] as String?,
      refreshInterval: (json['refreshInterval'] as num?)?.toInt(),
    );

Map<String, dynamic> _$$UpdateProjectRequestImplToJson(
        _$UpdateProjectRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'description': instance.description,
      'refreshInterval': instance.refreshInterval,
    };

_$ProjectBudgetRequestImpl _$$ProjectBudgetRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$ProjectBudgetRequestImpl(
      accountId: (json['accountId'] as num?)?.toInt(),
      budget: (json['budget'] as num).toDouble(),
      used: (json['used'] as num?)?.toDouble(),
    );

Map<String, dynamic> _$$ProjectBudgetRequestImplToJson(
        _$ProjectBudgetRequestImpl instance) =>
    <String, dynamic>{
      'accountId': instance.accountId,
      'budget': instance.budget,
      'used': instance.used,
    };

_$ProjectBudgetSummaryImpl _$$ProjectBudgetSummaryImplFromJson(
        Map<String, dynamic> json) =>
    _$ProjectBudgetSummaryImpl(
      budgets: (json['budgets'] as List<dynamic>)
          .map((e) => ProjectBudget.fromJson(e as Map<String, dynamic>))
          .toList(),
      totalBudget: (json['totalBudget'] as num).toDouble(),
      totalUsed: (json['totalUsed'] as num).toDouble(),
    );

Map<String, dynamic> _$$ProjectBudgetSummaryImplToJson(
        _$ProjectBudgetSummaryImpl instance) =>
    <String, dynamic>{
      'budgets': instance.budgets,
      'totalBudget': instance.totalBudget,
      'totalUsed': instance.totalUsed,
    };

_$ProjectListResponseImpl _$$ProjectListResponseImplFromJson(
        Map<String, dynamic> json) =>
    _$ProjectListResponseImpl(
      total: (json['total'] as num).toInt(),
      items: (json['items'] as List<dynamic>)
          .map((e) => Project.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$$ProjectListResponseImplToJson(
        _$ProjectListResponseImpl instance) =>
    <String, dynamic>{
      'total': instance.total,
      'items': instance.items,
    };
