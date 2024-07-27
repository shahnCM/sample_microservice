<?php

namespace App\Services;

class JwtService
{
    private $secret;

    public function __construct()
    {
        // Load the secret key from the environment file
        $this->secret = getenv('JWT_SECRET');
    }

    public function getClaimsFromToken(string $jwtToken): array
    {
        $parts = explode('.', $jwtToken);

        if (count($parts) !== 3) {
            throw new \Exception('Invalid JWT format');
        }

        $payload = $parts[1];

        return $this->base64UrlDecodeToArray($payload);
    }

    public function verifyAndGetClaimsFromToken(string $jwtToken): array
    {
        $parts = explode('.', $jwtToken);

        if (count($parts) !== 3) {
            throw new \Exception('Invalid JWT format');
        }

        list($header, $payload, $signature) = $parts;

        $decodedPayload = $this->base64UrlDecode($payload);

        $decodedSignature = $this->base64UrlDecode($signature);

        $expectedSignature = $this->generateSignature($header, $payload);

        if ($decodedSignature !== $expectedSignature) {
            throw new \Exception('Invalid token signature');
        }

        $claims = json_decode($decodedPayload, true);

        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new \Exception('Invalid token payload');
        }

        return $claims;
    }

    private function base64UrlDecodeToArray($input)
    {
        $remainder = strlen($input) % 4;
        if ($remainder) {
            $padLength = 4 - $remainder;
            $input .= str_repeat('=', $padLength);
        }
        return json_decode(base64_decode(strtr($input, '-_', '+/')), true);
    }

    private function base64UrlDecode($input)
    {
        $remainder = strlen($input) % 4;
        if ($remainder) {
            $padLength = 4 - $remainder;
            $input .= str_repeat('=', $padLength);
        }
        return base64_decode(strtr($input, '-_', '+/'));
    }

    private function generateSignature($header, $payload)
    {
        $headerAndPayload = $header . '.' . $payload;
        $hash = hash_hmac('sha256', $headerAndPayload, $this->secret, true);
        return $this->base64UrlEncode($hash);
    }

    private function base64UrlEncode($input)
    {
        return rtrim(strtr(base64_encode($input), '+/', '-_'), '=');
    }
}
