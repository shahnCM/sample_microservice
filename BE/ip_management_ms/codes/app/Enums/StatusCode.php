<?php
namespace App\Enums;

enum StatusCode: int
{
    case OK = 200;
    case CREATED = 201;
    case NO_CONTENT = 204;
    case BAD_REQUEST = 400;
    case UNAUTHORIZED = 401;
    case FORBIDDEN = 403;
    case NOT_FOUND = 404;
    case UNPROCESSABLE_ENTITY = 422;
    case INTERNAL_SERVER_ERROR = 500;
    case SERVICE_UNAVAILABLE = 503;

    public static function fromMessage(string $message): ?self
    {
        return match ($message) {
            'OK' => self::OK,
            'Created' => self::CREATED,
            'No Content' => self::NO_CONTENT,
            'Bad Request' => self::BAD_REQUEST,
            'Unauthorized' => self::UNAUTHORIZED,
            'Forbidden' => self::FORBIDDEN,
            'Not Found' => self::NOT_FOUND,
            'Unprocessable Entity' => self::UNPROCESSABLE_ENTITY,
            'Internal Server Error' => self::INTERNAL_SERVER_ERROR,
            'Service Unavailable' => self::SERVICE_UNAVAILABLE,
            default => null,
        };
    }

    public static function fromCode(int $code): ?self
    {
        return match ($code) {
            200 => self::OK,
            201 => self::CREATED,
            204 => self::NO_CONTENT,
            400 => self::BAD_REQUEST,
            401 => self::UNAUTHORIZED,
            403 => self::FORBIDDEN,
            404 => self::NOT_FOUND,
            422 => self::UNPROCESSABLE_ENTITY,
            500 => self::INTERNAL_SERVER_ERROR,
            503 => self::SERVICE_UNAVAILABLE,
            default => null,
        };
    }

    public function message(): string
    {
        return match ($this) {
            self::OK => 'OK',
            self::CREATED => 'Created',
            self::NO_CONTENT => 'No Content',
            self::BAD_REQUEST => 'Bad Request',
            self::UNAUTHORIZED => 'Unauthorized',
            self::FORBIDDEN => 'Forbidden',
            self::NOT_FOUND => 'Not Found',
            self::UNPROCESSABLE_ENTITY => 'Unprocessable Entity',
            self::INTERNAL_SERVER_ERROR => 'Internal Server Error',
            self::SERVICE_UNAVAILABLE => 'Service Unavailable',
        };
    }
}
