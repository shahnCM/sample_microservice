<?php

namespace App\Http\Middleware;

use App\Models\User;
use App\Services\JwtService;
use Auth;
use Closure;
use Symfony\Component\HttpKernel\Exception\UnauthorizedHttpException;

class ApiUserJwtAuth
{
    private $jwtService;

    public function __construct(JwtService $jwtService)
    {
        $this->jwtService = $jwtService;
    }

    public function handle($request, Closure $next)
    {
        $token = $request->header('Authorization');

        if (!$token) {
            throw new UnauthorizedHttpException('Invalid JWT format');
        }

        try {
            $claims = $this->jwtService->getClaimsFromToken($token);
            $user = new User();
            $user->id = $claims['user_id'];
            $user->sessionTokenTraceId = $claims['token_id'];
            Auth::setUser($user);

        } catch (\Exception $e) {
            throw new UnauthorizedHttpException('Invalid JWT format');
        }

        return $next($request);
    }
}
