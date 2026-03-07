// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'project.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

ProjectBudget _$ProjectBudgetFromJson(Map<String, dynamic> json) {
  return _ProjectBudget.fromJson(json);
}

/// @nodoc
mixin _$ProjectBudget {
  int get id => throw _privateConstructorUsedError;
  int get projectId => throw _privateConstructorUsedError;
  int? get accountId => throw _privateConstructorUsedError;
  double get budget => throw _privateConstructorUsedError;
  double get used => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this ProjectBudget to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of ProjectBudget
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $ProjectBudgetCopyWith<ProjectBudget> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ProjectBudgetCopyWith<$Res> {
  factory $ProjectBudgetCopyWith(
          ProjectBudget value, $Res Function(ProjectBudget) then) =
      _$ProjectBudgetCopyWithImpl<$Res, ProjectBudget>;
  @useResult
  $Res call(
      {int id,
      int projectId,
      int? accountId,
      double budget,
      double used,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class _$ProjectBudgetCopyWithImpl<$Res, $Val extends ProjectBudget>
    implements $ProjectBudgetCopyWith<$Res> {
  _$ProjectBudgetCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of ProjectBudget
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? projectId = null,
    Object? accountId = freezed,
    Object? budget = null,
    Object? used = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      projectId: null == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int,
      accountId: freezed == accountId
          ? _value.accountId
          : accountId // ignore: cast_nullable_to_non_nullable
              as int?,
      budget: null == budget
          ? _value.budget
          : budget // ignore: cast_nullable_to_non_nullable
              as double,
      used: null == used
          ? _value.used
          : used // ignore: cast_nullable_to_non_nullable
              as double,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$ProjectBudgetImplCopyWith<$Res>
    implements $ProjectBudgetCopyWith<$Res> {
  factory _$$ProjectBudgetImplCopyWith(
          _$ProjectBudgetImpl value, $Res Function(_$ProjectBudgetImpl) then) =
      __$$ProjectBudgetImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int projectId,
      int? accountId,
      double budget,
      double used,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class __$$ProjectBudgetImplCopyWithImpl<$Res>
    extends _$ProjectBudgetCopyWithImpl<$Res, _$ProjectBudgetImpl>
    implements _$$ProjectBudgetImplCopyWith<$Res> {
  __$$ProjectBudgetImplCopyWithImpl(
      _$ProjectBudgetImpl _value, $Res Function(_$ProjectBudgetImpl) _then)
      : super(_value, _then);

  /// Create a copy of ProjectBudget
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? projectId = null,
    Object? accountId = freezed,
    Object? budget = null,
    Object? used = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_$ProjectBudgetImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      projectId: null == projectId
          ? _value.projectId
          : projectId // ignore: cast_nullable_to_non_nullable
              as int,
      accountId: freezed == accountId
          ? _value.accountId
          : accountId // ignore: cast_nullable_to_non_nullable
              as int?,
      budget: null == budget
          ? _value.budget
          : budget // ignore: cast_nullable_to_non_nullable
              as double,
      used: null == used
          ? _value.used
          : used // ignore: cast_nullable_to_non_nullable
              as double,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$ProjectBudgetImpl implements _ProjectBudget {
  const _$ProjectBudgetImpl(
      {required this.id,
      required this.projectId,
      this.accountId,
      required this.budget,
      required this.used,
      required this.createdAt,
      required this.updatedAt});

  factory _$ProjectBudgetImpl.fromJson(Map<String, dynamic> json) =>
      _$$ProjectBudgetImplFromJson(json);

  @override
  final int id;
  @override
  final int projectId;
  @override
  final int? accountId;
  @override
  final double budget;
  @override
  final double used;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;

  @override
  String toString() {
    return 'ProjectBudget(id: $id, projectId: $projectId, accountId: $accountId, budget: $budget, used: $used, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$ProjectBudgetImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.projectId, projectId) ||
                other.projectId == projectId) &&
            (identical(other.accountId, accountId) ||
                other.accountId == accountId) &&
            (identical(other.budget, budget) || other.budget == budget) &&
            (identical(other.used, used) || other.used == used) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, id, projectId, accountId, budget,
      used, createdAt, updatedAt);

  /// Create a copy of ProjectBudget
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$ProjectBudgetImplCopyWith<_$ProjectBudgetImpl> get copyWith =>
      __$$ProjectBudgetImplCopyWithImpl<_$ProjectBudgetImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$ProjectBudgetImplToJson(
      this,
    );
  }
}

abstract class _ProjectBudget implements ProjectBudget {
  const factory _ProjectBudget(
      {required final int id,
      required final int projectId,
      final int? accountId,
      required final double budget,
      required final double used,
      required final DateTime createdAt,
      required final DateTime updatedAt}) = _$ProjectBudgetImpl;

  factory _ProjectBudget.fromJson(Map<String, dynamic> json) =
      _$ProjectBudgetImpl.fromJson;

  @override
  int get id;
  @override
  int get projectId;
  @override
  int? get accountId;
  @override
  double get budget;
  @override
  double get used;
  @override
  DateTime get createdAt;
  @override
  DateTime get updatedAt;

  /// Create a copy of ProjectBudget
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$ProjectBudgetImplCopyWith<_$ProjectBudgetImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

Project _$ProjectFromJson(Map<String, dynamic> json) {
  return _Project.fromJson(json);
}

/// @nodoc
mixin _$Project {
  int get id => throw _privateConstructorUsedError;
  int get userId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  int get refreshInterval => throw _privateConstructorUsedError;
  DateTime get lastRefresh => throw _privateConstructorUsedError;
  int get maxTaskId => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;
  List<ProjectBudget>? get budgets => throw _privateConstructorUsedError;

  /// Serializes this Project to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Project
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $ProjectCopyWith<Project> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ProjectCopyWith<$Res> {
  factory $ProjectCopyWith(Project value, $Res Function(Project) then) =
      _$ProjectCopyWithImpl<$Res, Project>;
  @useResult
  $Res call(
      {int id,
      int userId,
      String name,
      String? description,
      int refreshInterval,
      DateTime lastRefresh,
      int maxTaskId,
      DateTime createdAt,
      DateTime updatedAt,
      List<ProjectBudget>? budgets});
}

/// @nodoc
class _$ProjectCopyWithImpl<$Res, $Val extends Project>
    implements $ProjectCopyWith<$Res> {
  _$ProjectCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Project
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? name = null,
    Object? description = freezed,
    Object? refreshInterval = null,
    Object? lastRefresh = null,
    Object? maxTaskId = null,
    Object? createdAt = null,
    Object? updatedAt = null,
    Object? budgets = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      refreshInterval: null == refreshInterval
          ? _value.refreshInterval
          : refreshInterval // ignore: cast_nullable_to_non_nullable
              as int,
      lastRefresh: null == lastRefresh
          ? _value.lastRefresh
          : lastRefresh // ignore: cast_nullable_to_non_nullable
              as DateTime,
      maxTaskId: null == maxTaskId
          ? _value.maxTaskId
          : maxTaskId // ignore: cast_nullable_to_non_nullable
              as int,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      budgets: freezed == budgets
          ? _value.budgets
          : budgets // ignore: cast_nullable_to_non_nullable
              as List<ProjectBudget>?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$ProjectImplCopyWith<$Res> implements $ProjectCopyWith<$Res> {
  factory _$$ProjectImplCopyWith(
          _$ProjectImpl value, $Res Function(_$ProjectImpl) then) =
      __$$ProjectImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int userId,
      String name,
      String? description,
      int refreshInterval,
      DateTime lastRefresh,
      int maxTaskId,
      DateTime createdAt,
      DateTime updatedAt,
      List<ProjectBudget>? budgets});
}

/// @nodoc
class __$$ProjectImplCopyWithImpl<$Res>
    extends _$ProjectCopyWithImpl<$Res, _$ProjectImpl>
    implements _$$ProjectImplCopyWith<$Res> {
  __$$ProjectImplCopyWithImpl(
      _$ProjectImpl _value, $Res Function(_$ProjectImpl) _then)
      : super(_value, _then);

  /// Create a copy of Project
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? name = null,
    Object? description = freezed,
    Object? refreshInterval = null,
    Object? lastRefresh = null,
    Object? maxTaskId = null,
    Object? createdAt = null,
    Object? updatedAt = null,
    Object? budgets = freezed,
  }) {
    return _then(_$ProjectImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      refreshInterval: null == refreshInterval
          ? _value.refreshInterval
          : refreshInterval // ignore: cast_nullable_to_non_nullable
              as int,
      lastRefresh: null == lastRefresh
          ? _value.lastRefresh
          : lastRefresh // ignore: cast_nullable_to_non_nullable
              as DateTime,
      maxTaskId: null == maxTaskId
          ? _value.maxTaskId
          : maxTaskId // ignore: cast_nullable_to_non_nullable
              as int,
      createdAt: null == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      updatedAt: null == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as DateTime,
      budgets: freezed == budgets
          ? _value._budgets
          : budgets // ignore: cast_nullable_to_non_nullable
              as List<ProjectBudget>?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$ProjectImpl implements _Project {
  const _$ProjectImpl(
      {required this.id,
      required this.userId,
      required this.name,
      this.description,
      required this.refreshInterval,
      required this.lastRefresh,
      required this.maxTaskId,
      required this.createdAt,
      required this.updatedAt,
      final List<ProjectBudget>? budgets})
      : _budgets = budgets;

  factory _$ProjectImpl.fromJson(Map<String, dynamic> json) =>
      _$$ProjectImplFromJson(json);

  @override
  final int id;
  @override
  final int userId;
  @override
  final String name;
  @override
  final String? description;
  @override
  final int refreshInterval;
  @override
  final DateTime lastRefresh;
  @override
  final int maxTaskId;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;
  final List<ProjectBudget>? _budgets;
  @override
  List<ProjectBudget>? get budgets {
    final value = _budgets;
    if (value == null) return null;
    if (_budgets is EqualUnmodifiableListView) return _budgets;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(value);
  }

  @override
  String toString() {
    return 'Project(id: $id, userId: $userId, name: $name, description: $description, refreshInterval: $refreshInterval, lastRefresh: $lastRefresh, maxTaskId: $maxTaskId, createdAt: $createdAt, updatedAt: $updatedAt, budgets: $budgets)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$ProjectImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.userId, userId) || other.userId == userId) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.refreshInterval, refreshInterval) ||
                other.refreshInterval == refreshInterval) &&
            (identical(other.lastRefresh, lastRefresh) ||
                other.lastRefresh == lastRefresh) &&
            (identical(other.maxTaskId, maxTaskId) ||
                other.maxTaskId == maxTaskId) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt) &&
            const DeepCollectionEquality().equals(other._budgets, _budgets));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      userId,
      name,
      description,
      refreshInterval,
      lastRefresh,
      maxTaskId,
      createdAt,
      updatedAt,
      const DeepCollectionEquality().hash(_budgets));

  /// Create a copy of Project
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$ProjectImplCopyWith<_$ProjectImpl> get copyWith =>
      __$$ProjectImplCopyWithImpl<_$ProjectImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$ProjectImplToJson(
      this,
    );
  }
}

abstract class _Project implements Project {
  const factory _Project(
      {required final int id,
      required final int userId,
      required final String name,
      final String? description,
      required final int refreshInterval,
      required final DateTime lastRefresh,
      required final int maxTaskId,
      required final DateTime createdAt,
      required final DateTime updatedAt,
      final List<ProjectBudget>? budgets}) = _$ProjectImpl;

  factory _Project.fromJson(Map<String, dynamic> json) = _$ProjectImpl.fromJson;

  @override
  int get id;
  @override
  int get userId;
  @override
  String get name;
  @override
  String? get description;
  @override
  int get refreshInterval;
  @override
  DateTime get lastRefresh;
  @override
  int get maxTaskId;
  @override
  DateTime get createdAt;
  @override
  DateTime get updatedAt;
  @override
  List<ProjectBudget>? get budgets;

  /// Create a copy of Project
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$ProjectImplCopyWith<_$ProjectImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateProjectRequest _$CreateProjectRequestFromJson(Map<String, dynamic> json) {
  return _CreateProjectRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateProjectRequest {
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  int? get refreshInterval => throw _privateConstructorUsedError;

  /// Serializes this CreateProjectRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateProjectRequestCopyWith<CreateProjectRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateProjectRequestCopyWith<$Res> {
  factory $CreateProjectRequestCopyWith(CreateProjectRequest value,
          $Res Function(CreateProjectRequest) then) =
      _$CreateProjectRequestCopyWithImpl<$Res, CreateProjectRequest>;
  @useResult
  $Res call({String name, String? description, int? refreshInterval});
}

/// @nodoc
class _$CreateProjectRequestCopyWithImpl<$Res,
        $Val extends CreateProjectRequest>
    implements $CreateProjectRequestCopyWith<$Res> {
  _$CreateProjectRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? description = freezed,
    Object? refreshInterval = freezed,
  }) {
    return _then(_value.copyWith(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      refreshInterval: freezed == refreshInterval
          ? _value.refreshInterval
          : refreshInterval // ignore: cast_nullable_to_non_nullable
              as int?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateProjectRequestImplCopyWith<$Res>
    implements $CreateProjectRequestCopyWith<$Res> {
  factory _$$CreateProjectRequestImplCopyWith(_$CreateProjectRequestImpl value,
          $Res Function(_$CreateProjectRequestImpl) then) =
      __$$CreateProjectRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String name, String? description, int? refreshInterval});
}

/// @nodoc
class __$$CreateProjectRequestImplCopyWithImpl<$Res>
    extends _$CreateProjectRequestCopyWithImpl<$Res, _$CreateProjectRequestImpl>
    implements _$$CreateProjectRequestImplCopyWith<$Res> {
  __$$CreateProjectRequestImplCopyWithImpl(_$CreateProjectRequestImpl _value,
      $Res Function(_$CreateProjectRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? description = freezed,
    Object? refreshInterval = freezed,
  }) {
    return _then(_$CreateProjectRequestImpl(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      refreshInterval: freezed == refreshInterval
          ? _value.refreshInterval
          : refreshInterval // ignore: cast_nullable_to_non_nullable
              as int?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateProjectRequestImpl implements _CreateProjectRequest {
  const _$CreateProjectRequestImpl(
      {required this.name, this.description, this.refreshInterval});

  factory _$CreateProjectRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateProjectRequestImplFromJson(json);

  @override
  final String name;
  @override
  final String? description;
  @override
  final int? refreshInterval;

  @override
  String toString() {
    return 'CreateProjectRequest(name: $name, description: $description, refreshInterval: $refreshInterval)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateProjectRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.refreshInterval, refreshInterval) ||
                other.refreshInterval == refreshInterval));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, name, description, refreshInterval);

  /// Create a copy of CreateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateProjectRequestImplCopyWith<_$CreateProjectRequestImpl>
      get copyWith =>
          __$$CreateProjectRequestImplCopyWithImpl<_$CreateProjectRequestImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateProjectRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateProjectRequest implements CreateProjectRequest {
  const factory _CreateProjectRequest(
      {required final String name,
      final String? description,
      final int? refreshInterval}) = _$CreateProjectRequestImpl;

  factory _CreateProjectRequest.fromJson(Map<String, dynamic> json) =
      _$CreateProjectRequestImpl.fromJson;

  @override
  String get name;
  @override
  String? get description;
  @override
  int? get refreshInterval;

  /// Create a copy of CreateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateProjectRequestImplCopyWith<_$CreateProjectRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

UpdateProjectRequest _$UpdateProjectRequestFromJson(Map<String, dynamic> json) {
  return _UpdateProjectRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateProjectRequest {
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  int? get refreshInterval => throw _privateConstructorUsedError;

  /// Serializes this UpdateProjectRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateProjectRequestCopyWith<UpdateProjectRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateProjectRequestCopyWith<$Res> {
  factory $UpdateProjectRequestCopyWith(UpdateProjectRequest value,
          $Res Function(UpdateProjectRequest) then) =
      _$UpdateProjectRequestCopyWithImpl<$Res, UpdateProjectRequest>;
  @useResult
  $Res call({String name, String? description, int? refreshInterval});
}

/// @nodoc
class _$UpdateProjectRequestCopyWithImpl<$Res,
        $Val extends UpdateProjectRequest>
    implements $UpdateProjectRequestCopyWith<$Res> {
  _$UpdateProjectRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? description = freezed,
    Object? refreshInterval = freezed,
  }) {
    return _then(_value.copyWith(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      refreshInterval: freezed == refreshInterval
          ? _value.refreshInterval
          : refreshInterval // ignore: cast_nullable_to_non_nullable
              as int?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateProjectRequestImplCopyWith<$Res>
    implements $UpdateProjectRequestCopyWith<$Res> {
  factory _$$UpdateProjectRequestImplCopyWith(_$UpdateProjectRequestImpl value,
          $Res Function(_$UpdateProjectRequestImpl) then) =
      __$$UpdateProjectRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String name, String? description, int? refreshInterval});
}

/// @nodoc
class __$$UpdateProjectRequestImplCopyWithImpl<$Res>
    extends _$UpdateProjectRequestCopyWithImpl<$Res, _$UpdateProjectRequestImpl>
    implements _$$UpdateProjectRequestImplCopyWith<$Res> {
  __$$UpdateProjectRequestImplCopyWithImpl(_$UpdateProjectRequestImpl _value,
      $Res Function(_$UpdateProjectRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? description = freezed,
    Object? refreshInterval = freezed,
  }) {
    return _then(_$UpdateProjectRequestImpl(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      refreshInterval: freezed == refreshInterval
          ? _value.refreshInterval
          : refreshInterval // ignore: cast_nullable_to_non_nullable
              as int?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateProjectRequestImpl implements _UpdateProjectRequest {
  const _$UpdateProjectRequestImpl(
      {required this.name, this.description, this.refreshInterval});

  factory _$UpdateProjectRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateProjectRequestImplFromJson(json);

  @override
  final String name;
  @override
  final String? description;
  @override
  final int? refreshInterval;

  @override
  String toString() {
    return 'UpdateProjectRequest(name: $name, description: $description, refreshInterval: $refreshInterval)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateProjectRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.refreshInterval, refreshInterval) ||
                other.refreshInterval == refreshInterval));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, name, description, refreshInterval);

  /// Create a copy of UpdateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateProjectRequestImplCopyWith<_$UpdateProjectRequestImpl>
      get copyWith =>
          __$$UpdateProjectRequestImplCopyWithImpl<_$UpdateProjectRequestImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateProjectRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateProjectRequest implements UpdateProjectRequest {
  const factory _UpdateProjectRequest(
      {required final String name,
      final String? description,
      final int? refreshInterval}) = _$UpdateProjectRequestImpl;

  factory _UpdateProjectRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateProjectRequestImpl.fromJson;

  @override
  String get name;
  @override
  String? get description;
  @override
  int? get refreshInterval;

  /// Create a copy of UpdateProjectRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateProjectRequestImplCopyWith<_$UpdateProjectRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

ProjectBudgetRequest _$ProjectBudgetRequestFromJson(Map<String, dynamic> json) {
  return _ProjectBudgetRequest.fromJson(json);
}

/// @nodoc
mixin _$ProjectBudgetRequest {
  int? get accountId => throw _privateConstructorUsedError;
  double get budget => throw _privateConstructorUsedError;
  double? get used => throw _privateConstructorUsedError;

  /// Serializes this ProjectBudgetRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of ProjectBudgetRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $ProjectBudgetRequestCopyWith<ProjectBudgetRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ProjectBudgetRequestCopyWith<$Res> {
  factory $ProjectBudgetRequestCopyWith(ProjectBudgetRequest value,
          $Res Function(ProjectBudgetRequest) then) =
      _$ProjectBudgetRequestCopyWithImpl<$Res, ProjectBudgetRequest>;
  @useResult
  $Res call({int? accountId, double budget, double? used});
}

/// @nodoc
class _$ProjectBudgetRequestCopyWithImpl<$Res,
        $Val extends ProjectBudgetRequest>
    implements $ProjectBudgetRequestCopyWith<$Res> {
  _$ProjectBudgetRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of ProjectBudgetRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? accountId = freezed,
    Object? budget = null,
    Object? used = freezed,
  }) {
    return _then(_value.copyWith(
      accountId: freezed == accountId
          ? _value.accountId
          : accountId // ignore: cast_nullable_to_non_nullable
              as int?,
      budget: null == budget
          ? _value.budget
          : budget // ignore: cast_nullable_to_non_nullable
              as double,
      used: freezed == used
          ? _value.used
          : used // ignore: cast_nullable_to_non_nullable
              as double?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$ProjectBudgetRequestImplCopyWith<$Res>
    implements $ProjectBudgetRequestCopyWith<$Res> {
  factory _$$ProjectBudgetRequestImplCopyWith(_$ProjectBudgetRequestImpl value,
          $Res Function(_$ProjectBudgetRequestImpl) then) =
      __$$ProjectBudgetRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int? accountId, double budget, double? used});
}

/// @nodoc
class __$$ProjectBudgetRequestImplCopyWithImpl<$Res>
    extends _$ProjectBudgetRequestCopyWithImpl<$Res, _$ProjectBudgetRequestImpl>
    implements _$$ProjectBudgetRequestImplCopyWith<$Res> {
  __$$ProjectBudgetRequestImplCopyWithImpl(_$ProjectBudgetRequestImpl _value,
      $Res Function(_$ProjectBudgetRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of ProjectBudgetRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? accountId = freezed,
    Object? budget = null,
    Object? used = freezed,
  }) {
    return _then(_$ProjectBudgetRequestImpl(
      accountId: freezed == accountId
          ? _value.accountId
          : accountId // ignore: cast_nullable_to_non_nullable
              as int?,
      budget: null == budget
          ? _value.budget
          : budget // ignore: cast_nullable_to_non_nullable
              as double,
      used: freezed == used
          ? _value.used
          : used // ignore: cast_nullable_to_non_nullable
              as double?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$ProjectBudgetRequestImpl implements _ProjectBudgetRequest {
  const _$ProjectBudgetRequestImpl(
      {this.accountId, required this.budget, this.used});

  factory _$ProjectBudgetRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$ProjectBudgetRequestImplFromJson(json);

  @override
  final int? accountId;
  @override
  final double budget;
  @override
  final double? used;

  @override
  String toString() {
    return 'ProjectBudgetRequest(accountId: $accountId, budget: $budget, used: $used)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$ProjectBudgetRequestImpl &&
            (identical(other.accountId, accountId) ||
                other.accountId == accountId) &&
            (identical(other.budget, budget) || other.budget == budget) &&
            (identical(other.used, used) || other.used == used));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, accountId, budget, used);

  /// Create a copy of ProjectBudgetRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$ProjectBudgetRequestImplCopyWith<_$ProjectBudgetRequestImpl>
      get copyWith =>
          __$$ProjectBudgetRequestImplCopyWithImpl<_$ProjectBudgetRequestImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$ProjectBudgetRequestImplToJson(
      this,
    );
  }
}

abstract class _ProjectBudgetRequest implements ProjectBudgetRequest {
  const factory _ProjectBudgetRequest(
      {final int? accountId,
      required final double budget,
      final double? used}) = _$ProjectBudgetRequestImpl;

  factory _ProjectBudgetRequest.fromJson(Map<String, dynamic> json) =
      _$ProjectBudgetRequestImpl.fromJson;

  @override
  int? get accountId;
  @override
  double get budget;
  @override
  double? get used;

  /// Create a copy of ProjectBudgetRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$ProjectBudgetRequestImplCopyWith<_$ProjectBudgetRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

ProjectBudgetSummary _$ProjectBudgetSummaryFromJson(Map<String, dynamic> json) {
  return _ProjectBudgetSummary.fromJson(json);
}

/// @nodoc
mixin _$ProjectBudgetSummary {
  List<ProjectBudget> get budgets => throw _privateConstructorUsedError;
  double get totalBudget => throw _privateConstructorUsedError;
  double get totalUsed => throw _privateConstructorUsedError;

  /// Serializes this ProjectBudgetSummary to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of ProjectBudgetSummary
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $ProjectBudgetSummaryCopyWith<ProjectBudgetSummary> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ProjectBudgetSummaryCopyWith<$Res> {
  factory $ProjectBudgetSummaryCopyWith(ProjectBudgetSummary value,
          $Res Function(ProjectBudgetSummary) then) =
      _$ProjectBudgetSummaryCopyWithImpl<$Res, ProjectBudgetSummary>;
  @useResult
  $Res call(
      {List<ProjectBudget> budgets, double totalBudget, double totalUsed});
}

/// @nodoc
class _$ProjectBudgetSummaryCopyWithImpl<$Res,
        $Val extends ProjectBudgetSummary>
    implements $ProjectBudgetSummaryCopyWith<$Res> {
  _$ProjectBudgetSummaryCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of ProjectBudgetSummary
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? budgets = null,
    Object? totalBudget = null,
    Object? totalUsed = null,
  }) {
    return _then(_value.copyWith(
      budgets: null == budgets
          ? _value.budgets
          : budgets // ignore: cast_nullable_to_non_nullable
              as List<ProjectBudget>,
      totalBudget: null == totalBudget
          ? _value.totalBudget
          : totalBudget // ignore: cast_nullable_to_non_nullable
              as double,
      totalUsed: null == totalUsed
          ? _value.totalUsed
          : totalUsed // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$ProjectBudgetSummaryImplCopyWith<$Res>
    implements $ProjectBudgetSummaryCopyWith<$Res> {
  factory _$$ProjectBudgetSummaryImplCopyWith(_$ProjectBudgetSummaryImpl value,
          $Res Function(_$ProjectBudgetSummaryImpl) then) =
      __$$ProjectBudgetSummaryImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {List<ProjectBudget> budgets, double totalBudget, double totalUsed});
}

/// @nodoc
class __$$ProjectBudgetSummaryImplCopyWithImpl<$Res>
    extends _$ProjectBudgetSummaryCopyWithImpl<$Res, _$ProjectBudgetSummaryImpl>
    implements _$$ProjectBudgetSummaryImplCopyWith<$Res> {
  __$$ProjectBudgetSummaryImplCopyWithImpl(_$ProjectBudgetSummaryImpl _value,
      $Res Function(_$ProjectBudgetSummaryImpl) _then)
      : super(_value, _then);

  /// Create a copy of ProjectBudgetSummary
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? budgets = null,
    Object? totalBudget = null,
    Object? totalUsed = null,
  }) {
    return _then(_$ProjectBudgetSummaryImpl(
      budgets: null == budgets
          ? _value._budgets
          : budgets // ignore: cast_nullable_to_non_nullable
              as List<ProjectBudget>,
      totalBudget: null == totalBudget
          ? _value.totalBudget
          : totalBudget // ignore: cast_nullable_to_non_nullable
              as double,
      totalUsed: null == totalUsed
          ? _value.totalUsed
          : totalUsed // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$ProjectBudgetSummaryImpl implements _ProjectBudgetSummary {
  const _$ProjectBudgetSummaryImpl(
      {required final List<ProjectBudget> budgets,
      required this.totalBudget,
      required this.totalUsed})
      : _budgets = budgets;

  factory _$ProjectBudgetSummaryImpl.fromJson(Map<String, dynamic> json) =>
      _$$ProjectBudgetSummaryImplFromJson(json);

  final List<ProjectBudget> _budgets;
  @override
  List<ProjectBudget> get budgets {
    if (_budgets is EqualUnmodifiableListView) return _budgets;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_budgets);
  }

  @override
  final double totalBudget;
  @override
  final double totalUsed;

  @override
  String toString() {
    return 'ProjectBudgetSummary(budgets: $budgets, totalBudget: $totalBudget, totalUsed: $totalUsed)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$ProjectBudgetSummaryImpl &&
            const DeepCollectionEquality().equals(other._budgets, _budgets) &&
            (identical(other.totalBudget, totalBudget) ||
                other.totalBudget == totalBudget) &&
            (identical(other.totalUsed, totalUsed) ||
                other.totalUsed == totalUsed));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType,
      const DeepCollectionEquality().hash(_budgets), totalBudget, totalUsed);

  /// Create a copy of ProjectBudgetSummary
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$ProjectBudgetSummaryImplCopyWith<_$ProjectBudgetSummaryImpl>
      get copyWith =>
          __$$ProjectBudgetSummaryImplCopyWithImpl<_$ProjectBudgetSummaryImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$ProjectBudgetSummaryImplToJson(
      this,
    );
  }
}

abstract class _ProjectBudgetSummary implements ProjectBudgetSummary {
  const factory _ProjectBudgetSummary(
      {required final List<ProjectBudget> budgets,
      required final double totalBudget,
      required final double totalUsed}) = _$ProjectBudgetSummaryImpl;

  factory _ProjectBudgetSummary.fromJson(Map<String, dynamic> json) =
      _$ProjectBudgetSummaryImpl.fromJson;

  @override
  List<ProjectBudget> get budgets;
  @override
  double get totalBudget;
  @override
  double get totalUsed;

  /// Create a copy of ProjectBudgetSummary
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$ProjectBudgetSummaryImplCopyWith<_$ProjectBudgetSummaryImpl>
      get copyWith => throw _privateConstructorUsedError;
}

ProjectListResponse _$ProjectListResponseFromJson(Map<String, dynamic> json) {
  return _ProjectListResponse.fromJson(json);
}

/// @nodoc
mixin _$ProjectListResponse {
  int get total => throw _privateConstructorUsedError;
  List<Project> get items => throw _privateConstructorUsedError;

  /// Serializes this ProjectListResponse to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of ProjectListResponse
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $ProjectListResponseCopyWith<ProjectListResponse> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $ProjectListResponseCopyWith<$Res> {
  factory $ProjectListResponseCopyWith(
          ProjectListResponse value, $Res Function(ProjectListResponse) then) =
      _$ProjectListResponseCopyWithImpl<$Res, ProjectListResponse>;
  @useResult
  $Res call({int total, List<Project> items});
}

/// @nodoc
class _$ProjectListResponseCopyWithImpl<$Res, $Val extends ProjectListResponse>
    implements $ProjectListResponseCopyWith<$Res> {
  _$ProjectListResponseCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of ProjectListResponse
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? total = null,
    Object? items = null,
  }) {
    return _then(_value.copyWith(
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      items: null == items
          ? _value.items
          : items // ignore: cast_nullable_to_non_nullable
              as List<Project>,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$ProjectListResponseImplCopyWith<$Res>
    implements $ProjectListResponseCopyWith<$Res> {
  factory _$$ProjectListResponseImplCopyWith(_$ProjectListResponseImpl value,
          $Res Function(_$ProjectListResponseImpl) then) =
      __$$ProjectListResponseImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int total, List<Project> items});
}

/// @nodoc
class __$$ProjectListResponseImplCopyWithImpl<$Res>
    extends _$ProjectListResponseCopyWithImpl<$Res, _$ProjectListResponseImpl>
    implements _$$ProjectListResponseImplCopyWith<$Res> {
  __$$ProjectListResponseImplCopyWithImpl(_$ProjectListResponseImpl _value,
      $Res Function(_$ProjectListResponseImpl) _then)
      : super(_value, _then);

  /// Create a copy of ProjectListResponse
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? total = null,
    Object? items = null,
  }) {
    return _then(_$ProjectListResponseImpl(
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      items: null == items
          ? _value._items
          : items // ignore: cast_nullable_to_non_nullable
              as List<Project>,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$ProjectListResponseImpl implements _ProjectListResponse {
  const _$ProjectListResponseImpl(
      {required this.total, required final List<Project> items})
      : _items = items;

  factory _$ProjectListResponseImpl.fromJson(Map<String, dynamic> json) =>
      _$$ProjectListResponseImplFromJson(json);

  @override
  final int total;
  final List<Project> _items;
  @override
  List<Project> get items {
    if (_items is EqualUnmodifiableListView) return _items;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_items);
  }

  @override
  String toString() {
    return 'ProjectListResponse(total: $total, items: $items)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$ProjectListResponseImpl &&
            (identical(other.total, total) || other.total == total) &&
            const DeepCollectionEquality().equals(other._items, _items));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType, total, const DeepCollectionEquality().hash(_items));

  /// Create a copy of ProjectListResponse
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$ProjectListResponseImplCopyWith<_$ProjectListResponseImpl> get copyWith =>
      __$$ProjectListResponseImplCopyWithImpl<_$ProjectListResponseImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$ProjectListResponseImplToJson(
      this,
    );
  }
}

abstract class _ProjectListResponse implements ProjectListResponse {
  const factory _ProjectListResponse(
      {required final int total,
      required final List<Project> items}) = _$ProjectListResponseImpl;

  factory _ProjectListResponse.fromJson(Map<String, dynamic> json) =
      _$ProjectListResponseImpl.fromJson;

  @override
  int get total;
  @override
  List<Project> get items;

  /// Create a copy of ProjectListResponse
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$ProjectListResponseImplCopyWith<_$ProjectListResponseImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
