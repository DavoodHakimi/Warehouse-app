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
import { apiFetch, ApiError } from '@/lib/api'

type Form = {
  company_name: string
  full_name: string
  user_name: string
  password: string
  password_confirmation: string
  phone_number: string
  email: string
}

const EMPTY: Form = {
  company_name: '',
  full_name: '',
  user_name: '',
  password: '',
  password_confirmation: '',
  phone_number: '',
  email: '',
}

function validate(form: Form): string | null {
  if (form.company_name.length < 4 || form.company_name.length > 100)
    return 'نام شرکت باید بین ۴ تا ۱۰۰ کاراکتر باشد.'
  if (form.full_name.length < 4 || form.full_name.length > 100)
    return 'نام کامل باید بین ۴ تا ۱۰۰ کاراکتر باشد.'
  if (form.user_name.length < 5 || form.user_name.length > 16)
    return 'نام کاربری باید بین ۵ تا ۱۶ کاراکتر باشد.'
  if (form.password.length < 8) return 'رمز عبور باید حداقل ۸ کاراکتر باشد.'
  if (form.password !== form.password_confirmation)
    return 'رمز عبور و تکرار آن مطابقت ندارند.'
  if (!/^09\d{9}$/.test(form.phone_number))
    return 'شماره تلفن باید ۱۱ رقم و با ۰۹ شروع شود.'
  if (
    form.email.length < 8 ||
    form.email.length > 32 ||
    !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)
  )
    return 'ایمیل معتبر و بین ۸ تا ۳۲ کاراکتر وارد کنید.'
  return null
}

export default function SignupPage() {
  const router = useRouter()
  const [form, setForm] = useState<Form>(EMPTY)
  const [submitting, setSubmitting] = useState(false)

  function update(key: keyof Form, value: string) {
    setForm((f) => ({ ...f, [key]: value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const error = validate(form)
    if (error) {
      toast.error(error)
      return
    }
    setSubmitting(true)
    try {
      await apiFetch('/auth/signup', {
        method: 'POST',
        auth: false,
        body: form,
      })
      toast.success('شرکت با موفقیت ثبت شد. اکنون وارد شوید.')
      router.replace('/login')
    } catch (err) {
      const message =
        err instanceof ApiError ? err.message : 'خطا در ثبت‌نام.'
      toast.error(message)
    } finally {
      setSubmitting(false)
    }
  }

  const fields: {
    key: keyof Form
    label: string
    type?: string
    dir?: 'ltr' | 'rtl'
    placeholder?: string
  }[] = [
    { key: 'company_name', label: 'نام شرکت', placeholder: 'نام شرکت' },
    { key: 'full_name', label: 'نام کامل مدیر', placeholder: 'نام و نام خانوادگی' },
    { key: 'user_name', label: 'نام کاربری', dir: 'ltr', placeholder: 'username' },
    { key: 'email', label: 'ایمیل', type: 'email', dir: 'ltr', placeholder: 'mail@example.com' },
    { key: 'phone_number', label: 'شماره تلفن', dir: 'ltr', placeholder: '09123456789' },
    { key: 'password', label: 'رمز عبور', type: 'password', dir: 'ltr', placeholder: '••••••••' },
    { key: 'password_confirmation', label: 'تکرار رمز عبور', type: 'password', dir: 'ltr', placeholder: '••••••••' },
  ]

  return (
    <main className="flex min-h-screen items-center justify-center bg-secondary px-4 py-12">
      <div className="w-full max-w-lg">
        <div className="mb-6 flex flex-col items-center gap-3 text-center">
          <div className="flex size-14 items-center justify-center rounded-2xl bg-primary text-primary-foreground">
            <Warehouse className="size-7" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-foreground">ثبت‌نام شرکت جدید</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              با ثبت شرکت، یک حساب مدیر برای شما ساخته می‌شود
            </p>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>اطلاعات شرکت و مدیر</CardTitle>
            <CardDescription>تمامی فیلدها الزامی هستند.</CardDescription>
          </CardHeader>
          <form onSubmit={handleSubmit}>
            <CardContent className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              {fields.map((f) => (
                <div
                  key={f.key}
                  className={`flex flex-col gap-2 ${f.key === 'company_name' || f.key === 'full_name' ? 'sm:col-span-2' : ''}`}
                >
                  <Label htmlFor={f.key}>{f.label}</Label>
                  <Input
                    id={f.key}
                    type={f.type ?? 'text'}
                    dir={f.dir ?? 'rtl'}
                    className={f.dir === 'ltr' ? 'text-right' : ''}
                    placeholder={f.placeholder}
                    value={form[f.key]}
                    onChange={(e) => update(f.key, e.target.value)}
                  />
                </div>
              ))}
            </CardContent>
            <CardFooter className="mt-6 flex flex-col gap-4">
              <Button type="submit" className="w-full" disabled={submitting}>
                {submitting && <Loader2 className="size-4 animate-spin" />}
                ثبت‌نام
              </Button>
              <p className="text-center text-sm text-muted-foreground">
                قبلاً ثبت‌نام کرده‌اید؟{' '}
                <Link
                  href="/login"
                  className="font-medium text-primary hover:underline"
                >
                  ورود به سامانه
                </Link>
              </p>
            </CardFooter>
          </form>
        </Card>
      </div>
    </main>
  )
}
