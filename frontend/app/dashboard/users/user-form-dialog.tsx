'use client'

import { useEffect, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { apiFetch, ApiError } from '@/lib/api'
import { ROLE_OPTIONS } from '@/lib/constants'
import type { User } from '@/lib/types'

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  user?: User | null
  onSaved: () => void
}

const empty = {
  full_name: '',
  user_name: '',
  user_type_id: 2,
  password: '',
  password_confirmation: '',
  phone_number: '',
  email: '',
}

export function UserFormDialog({ open, onOpenChange, user, onSaved }: Props) {
  const isEdit = Boolean(user)
  const [form, setForm] = useState(empty)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (open) {
      setForm(
        user
          ? {
              full_name: user.full_name,
              user_name: user.user_name,
              user_type_id: user.user_type_id,
              password: '',
              password_confirmation: '',
              phone_number: user.phone_number,
              email: user.email,
            }
          : empty,
      )
    }
  }, [open, user])

  function set<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }))
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()

    if (!isEdit && form.password !== form.password_confirmation) {
      toast.error('رمز عبور و تکرار آن یکسان نیستند.')
      return
    }

    setSaving(true)
    try {
      if (isEdit && user) {
        await apiFetch(`/users/${user.id}`, {
          method: 'PATCH',
          body: {
            id: user.id,
            full_name: form.full_name,
            user_name: form.user_name,
            user_type_id: Number(form.user_type_id),
            phone_number: form.phone_number,
            email: form.email,
          },
        })
        toast.success('کاربر با موفقیت ویرایش شد.')
      } else {
        await apiFetch('/users/', {
          method: 'POST',
          body: {
            full_name: form.full_name,
            user_name: form.user_name,
            user_type_id: Number(form.user_type_id),
            password: form.password,
            password_confirmation: form.password_confirmation,
            phone_number: form.phone_number,
            email: form.email,
          },
        })
        toast.success('کاربر جدید با موفقیت ایجاد شد.')
      }
      onSaved()
      onOpenChange(false)
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در ذخیره‌ی کاربر.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{isEdit ? 'ویرایش کاربر' : 'افزودن کاربر'}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? 'اطلاعات کاربر را ویرایش کنید.'
              : 'اطلاعات کاربر جدید را وارد کنید.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="full_name">نام و نام خانوادگی</Label>
            <Input
              id="full_name"
              value={form.full_name}
              onChange={(e) => set('full_name', e.target.value)}
              placeholder="مثلاً علی رضایی"
              required
              minLength={4}
            />
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="flex flex-col gap-2">
              <Label htmlFor="user_name">نام کاربری</Label>
              <Input
                id="user_name"
                dir="ltr"
                className="text-left"
                value={form.user_name}
                onChange={(e) => set('user_name', e.target.value)}
                placeholder="username"
                required
                minLength={5}
                maxLength={16}
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="role">نقش</Label>
              <Select
                value={String(form.user_type_id)}
                onValueChange={(v) => set('user_type_id', Number(v))}
              >
                <SelectTrigger id="role">
                  <SelectValue placeholder="انتخاب نقش" />
                </SelectTrigger>
                <SelectContent>
                  {ROLE_OPTIONS.map((r) => (
                    <SelectItem key={r.value} value={String(r.value)}>
                      {r.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {!isEdit && (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="flex flex-col gap-2">
                <Label htmlFor="password">رمز عبور</Label>
                <Input
                  id="password"
                  type="password"
                  dir="ltr"
                  className="text-left"
                  value={form.password}
                  onChange={(e) => set('password', e.target.value)}
                  placeholder="حداقل ۸ کاراکتر"
                  required
                  minLength={8}
                />
              </div>
              <div className="flex flex-col gap-2">
                <Label htmlFor="password_confirmation">تکرار رمز عبور</Label>
                <Input
                  id="password_confirmation"
                  type="password"
                  dir="ltr"
                  className="text-left"
                  value={form.password_confirmation}
                  onChange={(e) => set('password_confirmation', e.target.value)}
                  required
                  minLength={8}
                />
              </div>
            </div>
          )}

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="flex flex-col gap-2">
              <Label htmlFor="phone_number">شماره تماس</Label>
              <Input
                id="phone_number"
                dir="ltr"
                className="text-left"
                value={form.phone_number}
                onChange={(e) => set('phone_number', e.target.value)}
                placeholder="09xxxxxxxxx"
                required
                pattern="09[0-9]{9}"
                title="شماره باید ۱۱ رقم و با ۰۹ شروع شود"
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="email">ایمیل</Label>
              <Input
                id="email"
                type="email"
                dir="ltr"
                className="text-left"
                value={form.email}
                onChange={(e) => set('email', e.target.value)}
                placeholder="name@example.com"
                required
              />
            </div>
          </div>

          <DialogFooter className="mt-2 gap-2 sm:gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              انصراف
            </Button>
            <Button type="submit" disabled={saving}>
              {saving && <Loader2 className="size-4 animate-spin" />}
              {isEdit ? 'ذخیره تغییرات' : 'ایجاد کاربر'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
