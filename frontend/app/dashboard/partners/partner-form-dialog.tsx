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
import { PARTNER_TYPE_OPTIONS } from '@/lib/constants'
import type { Partner } from '@/lib/types'

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  partner?: Partner | null
  onSaved: () => void
}

const empty = {
  name: '',
  business_partner_type_id: 1,
  phone_number: '',
  email: '',
  contact_name: '',
  contact_phone_number: '',
}

// نگاشت نام انگلیسی نوع شریک به شناسه
const TYPE_NAME_TO_ID: Record<string, number> = {
  Supplier: 1,
  Customer: 2,
  Both: 3,
}

export function PartnerFormDialog({
  open,
  onOpenChange,
  partner,
  onSaved,
}: Props) {
  const isEdit = Boolean(partner)
  const [form, setForm] = useState(empty)
  const [saving, setSaving] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (open) {
      setErrors({})
      setForm(
        partner
          ? {
              name: partner.name,
              business_partner_type_id:
                TYPE_NAME_TO_ID[partner.business_partner_type] ?? 1,
              phone_number: partner.phone_number,
              email: partner.email,
              contact_name: partner.contact_name,
              contact_phone_number: partner.contact_phone_number,
            }
          : empty,
      )
    }
  }, [open, partner])

  function set<K extends keyof typeof form>(key: K, value: (typeof form)[K]) {
    setForm((f) => ({ ...f, [key]: value }))
    setErrors((e) => ({ ...e, [key]: '' }))
  }

  function validate(): boolean {
    const errs: Record<string, string> = {}
    const phoneRegex = /^09\d{9}$/
    if (!phoneRegex.test(form.phone_number))
      errs.phone_number = 'شماره تماس باید ۱۱ رقم و با ۰۹ شروع شود'
    if (
      form.contact_phone_number &&
      !/^09\d{9}$/.test(form.contact_phone_number)
    )
      errs.contact_phone_number =
        'شماره شخص رابط باید ۱۱ رقم و با ۰۹ شروع شود'
    setErrors(errs)
    return Object.keys(errs).length === 0
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!validate()) return
    setSaving(true)
    try {
      const body = {
        name: form.name,
        business_partner_type_id: Number(form.business_partner_type_id),
        phone_number: form.phone_number,
        email: form.email,
        contact_name: form.contact_name,
        contact_phone_number: form.contact_phone_number,
      }
      if (isEdit && partner) {
        await apiFetch(`/partners/${partner.id}`, {
          method: 'PATCH',
          body: { id: partner.id, ...body },
        })
        toast.success('شریک تجاری ویرایش شد.')
      } else {
        await apiFetch('/partners/', { method: 'POST', body })
        toast.success('شریک تجاری جدید ایجاد شد.')
      }
      onSaved()
      onOpenChange(false)
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در ذخیره.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>
            {isEdit ? 'ویرایش شریک تجاری' : 'افزودن شریک تجاری'}
          </DialogTitle>
          <DialogDescription>
            اطلاعات شریک تجاری (تأمین‌کننده یا خریدار) را وارد کنید.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="flex flex-col gap-2">
              <Label htmlFor="name">نام شرکت</Label>
              <Input
                id="name"
                value={form.name}
                onChange={(e) => set('name', e.target.value)}
                placeholder="مثلاً شرکت آلفا"
                required
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="type">نوع شریک</Label>
              <Select
                value={String(form.business_partner_type_id)}
                onValueChange={(v) =>
                  set('business_partner_type_id', Number(v))
                }
              >
                <SelectTrigger id="type">
                  <SelectValue placeholder="انتخاب نوع" />
                </SelectTrigger>
                <SelectContent>
                  {PARTNER_TYPE_OPTIONS.map((t) => (
                    <SelectItem key={t.value} value={String(t.value)}>
                      {t.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

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
              />
              {errors.phone_number && (
                <span className="text-xs text-destructive">
                  {errors.phone_number}
                </span>
              )}
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
                placeholder="info@example.com"
                required
              />
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="flex flex-col gap-2">
              <Label htmlFor="contact_name">نام شخص رابط</Label>
              <Input
                id="contact_name"
                value={form.contact_name}
                onChange={(e) => set('contact_name', e.target.value)}
                placeholder="مثلاً علی رضایی"
                required
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="contact_phone_number">شماره شخص رابط</Label>
              <Input
                id="contact_phone_number"
                dir="ltr"
                className="text-left"
                value={form.contact_phone_number}
                onChange={(e) => set('contact_phone_number', e.target.value)}
                placeholder="09xxxxxxxxx"
                required
              />
              {errors.contact_phone_number && (
                <span className="text-xs text-destructive">
                  {errors.contact_phone_number}
                </span>
              )}
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
              {isEdit ? 'ذخیره تغییرات' : 'ایجاد شریک'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
