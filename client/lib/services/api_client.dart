import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/models.dart';

class ApiClient {
  late final Dio _dio;
  String? _accessToken;
  String? _refreshToken;
  final void Function()? onTokenExpired;

  ApiClient({
    String baseUrl = 'http://localhost:8080/api/v1',
    this.onTokenExpired,
  }) {
    _dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: const Duration(seconds: 30),
        receiveTimeout: const Duration(seconds: 30),
        headers: {
          'Content-Type': 'application/json',
        },
      ),
    );

    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) {
          if (_accessToken != null) {
            options.headers['Authorization'] = 'Bearer $_accessToken';
          }
          return handler.next(options);
        },
        onError: (error, handler) async {
          if (error.response?.statusCode == 401 && _refreshToken != null) {
            try {
              final response = await _dio.post(
                '/auth/refresh',
                data: {'refresh_token': _refreshToken},
              );
              final authResponse = AuthResponse.fromJson(response.data['data']);
              setTokens(authResponse.accessToken, authResponse.refreshToken);
              error.requestOptions.headers['Authorization'] =
                  'Bearer $_accessToken';
              final retryResponse = await _dio.fetch(error.requestOptions);
              return handler.resolve(retryResponse);
            } catch (e) {
              onTokenExpired?.call();
              return handler.reject(error);
            }
          }
          return handler.next(error);
        },
      ),
    );
  }

  void setTokens(String accessToken, String refreshToken) {
    _accessToken = accessToken;
    _refreshToken = refreshToken;
  }

  void clearTokens() {
    _accessToken = null;
    _refreshToken = null;
  }

  Dio get dio => _dio;
}

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient(
    onTokenExpired: () {
      ref.read(authStateProvider.notifier).logout();
    },
  );
});

class ApiException implements Exception {
  final int code;
  final String message;

  ApiException(this.code, this.message);

  @override
  String toString() => 'ApiException: $code - $message';
}
