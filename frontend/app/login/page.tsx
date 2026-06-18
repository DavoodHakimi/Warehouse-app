'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { Loader2, Warehouse } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { useAuth } from '@/lib/auth'
import { ApiError } from '@/lib/api'

export default function LoginPage() {
  const { login } = useAuth()
  const router = useRouter()
  const [userName, setUserName] = useState('')
  const [password, setPassword] = useState('')
  const [submitting, setSubmitting] = useState(false)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!userName || !password) {
      toast.error('نام کاربری و رمز عبور را وارد کنید.')
      return
    }
    setSubmitting(true)
    try {
      await login(userName, password)
      toast.success('با موفقیت وارد شدید.')
      router.replace('/dashboard')
    } catch (err) {
      const message =
        err instanceof ApiError
          ? err.status === 404
            ? 'کاربری با این نام یافت نشد.'
            : err.status === 401
              ? 'رمز عبور نادرست است.'
              : err.message
          : 'خطا در ورود.'
      toast.error(message)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-secondary px-4 py-12">
      <div className="w-full max-w-md">
        <div className="mb-6 flex flex-col items-center gap-3 text-center">
          <div className="flex size-14 items-center justify-center rounded-2xl bg-primary text-primary-foreground">
            <Warehouse className="size-7" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-foreground">سامانه مدیریت انبار</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              برای ادامه وارد حساب کاربری خود شوید
            </p>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>ورود به سامانه</CardTitle>
            <CardDescription>نام کاربری و رمز عبور خود را وارد کنید.</CardDescription>
          </CardHeader>
          <form onSubmit={handleSubmit}>
            <CardContent className="flex flex-col gap-4">
              <div className="flex flex-col gap-2">
                <Label htmlFor="userName">نام کاربری</Label>
                <Input
                  id="userName"
                  value={userName}
                  onChange={(e) => setUserName(e.target.value)}
                  placeholder="نام کاربری"
                  autoComplete="username"
                  dir="ltr"
                  className="text-right"
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="password">رمز عبور</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="••••••••"
                  autoComplete="current-password"
                  dir="ltr"
                  className="text-right"
                />
              </div>
            </CardContent>
            <CardFooter className="mt-6 flex flex-col gap-4">
              <Button type="submit" className="w-full" disabled={submitting}>
                {submitting && <Loader2 className="size-4 animate-spin" />}
                ورود
              </Button>
              <p className="text-center text-sm text-muted-foreground">
                حساب کاربری ندارید؟{' '}
                <Link
                  href="/signup"
                  className="font-medium text-primary hover:underline"
                >
                  ثبت‌نام شرکت جدید
                </Link>
              </p>
            </CardFooter>
          </form>
        </Card>
      </div>
    </main>
  )
}
