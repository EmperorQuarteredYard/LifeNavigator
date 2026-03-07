// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'account.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Account _$AccountFromJson(Map<String, dynamic> json) {
  return _Account.fromJson(json);
}

/// @nodoc
mixin _$Account {
  int get id => throw _privateConstructorUsedError;
  int get userId => throw _privateConstructorUsedError;
  String get type => throw _privateConstructorUsedError;
  double get balance => throw _privateConstructorUsedError;
  int get version => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  DateTime get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this Account to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Account
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $AccountCopyWith<Account> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $AccountCopyWith<$Res> {
  factory $AccountCopyWith(Account value, $Res Function(Account) then) =
      _$AccountCopyWithImpl<$Res, Account>;
  @useResult
  $Res call(
      {int id,
      int userId,
      String type,
      double balance,
      int version,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class _$AccountCopyWithImpl<$Res, $Val extends Account>
    implements $AccountCopyWith<$Res> {
  _$AccountCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Account
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? type = null,
    Object? balance = null,
    Object? version = null,
    Object? createdAt = null,
    Object? updatedAt = null,
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
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      balance: null == balance
          ? _value.balance
          : balance // ignore: cast_nullable_to_non_nullable
              as double,
      version: null == version
          ? _value.version
          : version // ignore: cast_nullable_to_non_nullable
              as int,
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
abstract class _$$AccountImplCopyWith<$Res> implements $AccountCopyWith<$Res> {
  factory _$$AccountImplCopyWith(
          _$AccountImpl value, $Res Function(_$AccountImpl) then) =
      __$$AccountImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int userId,
      String type,
      double balance,
      int version,
      DateTime createdAt,
      DateTime updatedAt});
}

/// @nodoc
class __$$AccountImplCopyWithImpl<$Res>
    extends _$AccountCopyWithImpl<$Res, _$AccountImpl>
    implements _$$AccountImplCopyWith<$Res> {
  __$$AccountImplCopyWithImpl(
      _$AccountImpl _value, $Res Function(_$AccountImpl) _then)
      : super(_value, _then);

  /// Create a copy of Account
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? userId = null,
    Object? type = null,
    Object? balance = null,
    Object? version = null,
    Object? createdAt = null,
    Object? updatedAt = null,
  }) {
    return _then(_$AccountImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      userId: null == userId
          ? _value.userId
          : userId // ignore: cast_nullable_to_non_nullable
              as int,
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      balance: null == balance
          ? _value.balance
          : balance // ignore: cast_nullable_to_non_nullable
              as double,
      version: null == version
          ? _value.version
          : version // ignore: cast_nullable_to_non_nullable
              as int,
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
class _$AccountImpl implements _Account {
  const _$AccountImpl(
      {required this.id,
      required this.userId,
      required this.type,
      required this.balance,
      required this.version,
      required this.createdAt,
      required this.updatedAt});

  factory _$AccountImpl.fromJson(Map<String, dynamic> json) =>
      _$$AccountImplFromJson(json);

  @override
  final int id;
  @override
  final int userId;
  @override
  final String type;
  @override
  final double balance;
  @override
  final int version;
  @override
  final DateTime createdAt;
  @override
  final DateTime updatedAt;

  @override
  String toString() {
    return 'Account(id: $id, userId: $userId, type: $type, balance: $balance, version: $version, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$AccountImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.userId, userId) || other.userId == userId) &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.balance, balance) || other.balance == balance) &&
            (identical(other.version, version) || other.version == version) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType, id, userId, type, balance, version, createdAt, updatedAt);

  /// Create a copy of Account
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$AccountImplCopyWith<_$AccountImpl> get copyWith =>
      __$$AccountImplCopyWithImpl<_$AccountImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$AccountImplToJson(
      this,
    );
  }
}

abstract class _Account implements Account {
  const factory _Account(
      {required final int id,
      required final int userId,
      required final String type,
      required final double balance,
      required final int version,
      required final DateTime createdAt,
      required final DateTime updatedAt}) = _$AccountImpl;

  factory _Account.fromJson(Map<String, dynamic> json) = _$AccountImpl.fromJson;

  @override
  int get id;
  @override
  int get userId;
  @override
  String get type;
  @override
  double get balance;
  @override
  int get version;
  @override
  DateTime get createdAt;
  @override
  DateTime get updatedAt;

  /// Create a copy of Account
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$AccountImplCopyWith<_$AccountImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateAccountRequest _$CreateAccountRequestFromJson(Map<String, dynamic> json) {
  return _CreateAccountRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateAccountRequest {
  String get type => throw _privateConstructorUsedError;
  double? get balance => throw _privateConstructorUsedError;

  /// Serializes this CreateAccountRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateAccountRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateAccountRequestCopyWith<CreateAccountRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateAccountRequestCopyWith<$Res> {
  factory $CreateAccountRequestCopyWith(CreateAccountRequest value,
          $Res Function(CreateAccountRequest) then) =
      _$CreateAccountRequestCopyWithImpl<$Res, CreateAccountRequest>;
  @useResult
  $Res call({String type, double? balance});
}

/// @nodoc
class _$CreateAccountRequestCopyWithImpl<$Res,
        $Val extends CreateAccountRequest>
    implements $CreateAccountRequestCopyWith<$Res> {
  _$CreateAccountRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateAccountRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? balance = freezed,
  }) {
    return _then(_value.copyWith(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      balance: freezed == balance
          ? _value.balance
          : balance // ignore: cast_nullable_to_non_nullable
              as double?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateAccountRequestImplCopyWith<$Res>
    implements $CreateAccountRequestCopyWith<$Res> {
  factory _$$CreateAccountRequestImplCopyWith(_$CreateAccountRequestImpl value,
          $Res Function(_$CreateAccountRequestImpl) then) =
      __$$CreateAccountRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String type, double? balance});
}

/// @nodoc
class __$$CreateAccountRequestImplCopyWithImpl<$Res>
    extends _$CreateAccountRequestCopyWithImpl<$Res, _$CreateAccountRequestImpl>
    implements _$$CreateAccountRequestImplCopyWith<$Res> {
  __$$CreateAccountRequestImplCopyWithImpl(_$CreateAccountRequestImpl _value,
      $Res Function(_$CreateAccountRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateAccountRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? type = null,
    Object? balance = freezed,
  }) {
    return _then(_$CreateAccountRequestImpl(
      type: null == type
          ? _value.type
          : type // ignore: cast_nullable_to_non_nullable
              as String,
      balance: freezed == balance
          ? _value.balance
          : balance // ignore: cast_nullable_to_non_nullable
              as double?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateAccountRequestImpl implements _CreateAccountRequest {
  const _$CreateAccountRequestImpl({required this.type, this.balance});

  factory _$CreateAccountRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateAccountRequestImplFromJson(json);

  @override
  final String type;
  @override
  final double? balance;

  @override
  String toString() {
    return 'CreateAccountRequest(type: $type, balance: $balance)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateAccountRequestImpl &&
            (identical(other.type, type) || other.type == type) &&
            (identical(other.balance, balance) || other.balance == balance));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, type, balance);

  /// Create a copy of CreateAccountRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateAccountRequestImplCopyWith<_$CreateAccountRequestImpl>
      get copyWith =>
          __$$CreateAccountRequestImplCopyWithImpl<_$CreateAccountRequestImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateAccountRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateAccountRequest implements CreateAccountRequest {
  const factory _CreateAccountRequest(
      {required final String type,
      final double? balance}) = _$CreateAccountRequestImpl;

  factory _CreateAccountRequest.fromJson(Map<String, dynamic> json) =
      _$CreateAccountRequestImpl.fromJson;

  @override
  String get type;
  @override
  double? get balance;

  /// Create a copy of CreateAccountRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateAccountRequestImplCopyWith<_$CreateAccountRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

AdjustBalanceRequest _$AdjustBalanceRequestFromJson(Map<String, dynamic> json) {
  return _AdjustBalanceRequest.fromJson(json);
}

/// @nodoc
mixin _$AdjustBalanceRequest {
  double get amount => throw _privateConstructorUsedError;

  /// Serializes this AdjustBalanceRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of AdjustBalanceRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $AdjustBalanceRequestCopyWith<AdjustBalanceRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $AdjustBalanceRequestCopyWith<$Res> {
  factory $AdjustBalanceRequestCopyWith(AdjustBalanceRequest value,
          $Res Function(AdjustBalanceRequest) then) =
      _$AdjustBalanceRequestCopyWithImpl<$Res, AdjustBalanceRequest>;
  @useResult
  $Res call({double amount});
}

/// @nodoc
class _$AdjustBalanceRequestCopyWithImpl<$Res,
        $Val extends AdjustBalanceRequest>
    implements $AdjustBalanceRequestCopyWith<$Res> {
  _$AdjustBalanceRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of AdjustBalanceRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? amount = null,
  }) {
    return _then(_value.copyWith(
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$AdjustBalanceRequestImplCopyWith<$Res>
    implements $AdjustBalanceRequestCopyWith<$Res> {
  factory _$$AdjustBalanceRequestImplCopyWith(_$AdjustBalanceRequestImpl value,
          $Res Function(_$AdjustBalanceRequestImpl) then) =
      __$$AdjustBalanceRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({double amount});
}

/// @nodoc
class __$$AdjustBalanceRequestImplCopyWithImpl<$Res>
    extends _$AdjustBalanceRequestCopyWithImpl<$Res, _$AdjustBalanceRequestImpl>
    implements _$$AdjustBalanceRequestImplCopyWith<$Res> {
  __$$AdjustBalanceRequestImplCopyWithImpl(_$AdjustBalanceRequestImpl _value,
      $Res Function(_$AdjustBalanceRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of AdjustBalanceRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? amount = null,
  }) {
    return _then(_$AdjustBalanceRequestImpl(
      amount: null == amount
          ? _value.amount
          : amount // ignore: cast_nullable_to_non_nullable
              as double,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$AdjustBalanceRequestImpl implements _AdjustBalanceRequest {
  const _$AdjustBalanceRequestImpl({required this.amount});

  factory _$AdjustBalanceRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$AdjustBalanceRequestImplFromJson(json);

  @override
  final double amount;

  @override
  String toString() {
    return 'AdjustBalanceRequest(amount: $amount)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$AdjustBalanceRequestImpl &&
            (identical(other.amount, amount) || other.amount == amount));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, amount);

  /// Create a copy of AdjustBalanceRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$AdjustBalanceRequestImplCopyWith<_$AdjustBalanceRequestImpl>
      get copyWith =>
          __$$AdjustBalanceRequestImplCopyWithImpl<_$AdjustBalanceRequestImpl>(
              this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$AdjustBalanceRequestImplToJson(
      this,
    );
  }
}

abstract class _AdjustBalanceRequest implements AdjustBalanceRequest {
  const factory _AdjustBalanceRequest({required final double amount}) =
      _$AdjustBalanceRequestImpl;

  factory _AdjustBalanceRequest.fromJson(Map<String, dynamic> json) =
      _$AdjustBalanceRequestImpl.fromJson;

  @override
  double get amount;

  /// Create a copy of AdjustBalanceRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$AdjustBalanceRequestImplCopyWith<_$AdjustBalanceRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}
