'use client'

import { useEffect, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { apiFetch, ApiError } from '@/lib/api'
import type { Product } from '@/lib/types'

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  product?: Product | null
  onSaved: () => void
}

export function ProductFormDialog({
  open,
  onOpenChange,
  product,
  onSaved,
}: Props) {
  const isEdit = Boolean(product)
  const [name, setName] = useState('')
  const [price, setPrice] = useState('')
  const [isFrozen, setIsFrozen] = useState(false)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (open) {
      setName(product?.name ?? '')
      setPrice(product ? String(product.default_price) : '')
      setIsFrozen(product?.is_frozen ?? false)
    }
  }, [open, product])

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      if (isEdit && product) {
        await apiFetch(`/products/${product.product_number}/`, {
          method: 'PATCH',
          body: {
            id: product.id,
            name,
            product_number: product.product_number,
            is_frozen: isFrozen,
            default_price: Number(price),
          },
        })
        toast.success('محصول ویرایش شد.')
      } else {
        await apiFetch('/products/', {
          method: 'POST',
          body: {
            name,
            is_frozen: isFrozen,
            default_price: Number(price),
          },
        })
        toast.success('محصول جدید ایجاد شد.')
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
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{isEdit ? 'ویرایش محصول' : 'افزودن محصول'}</DialogTitle>
          <DialogDescription>
            {isEdit
              ? 'اطلاعات محصول را ویرایش کنید.'
              : 'شماره محصول به‌صورت خودکار ساخته می‌شود.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="name">نام محصول</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="مثلاً پیچ ۱۰ میلی‌متری"
              required
              minLength={4}
            />
          </div>

          {isEdit && product && (
            <div className="flex flex-col gap-2">
              <Label>شماره محصول</Label>
              <Input
                value={product.product_number}
                dir="ltr"
                className="text-left"
                disabled
              />
            </div>
          )}

          <div className="flex flex-col gap-2">
            <Label htmlFor="price">قیمت پایه (ریال)</Label>
            <Input
              id="price"
              type="number"
              min={0}
              dir="ltr"
              className="text-left"
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              placeholder="150000"
              required
            />
          </div>

          <div className="flex items-center justify-between rounded-lg border border-border p-3">
            <div className="flex flex-col">
              <Label htmlFor="frozen">قفل کردن محصول</Label>
              <span className="text-xs text-muted-foreground">
                محصول قفل شده در سفارش‌های جدید قابل استفاده نیست.
              </span>
            </div>
            <Switch id="frozen" checked={isFrozen} onCheckedChange={setIsFrozen}/>
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
              {isEdit ? 'ذخیره تغییرات' : 'ایجاد محصول'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
