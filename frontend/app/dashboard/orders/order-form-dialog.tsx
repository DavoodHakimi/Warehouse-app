'use client'

import { useEffect, useState } from 'react'
import { Loader2, Plus, Trash2 } from 'lucide-react'
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
import { CURRENCY_OPTIONS, ORDER_TYPES } from '@/lib/constants'
import { formatPrice } from '@/lib/utils'
import type { Partner, Product } from '@/lib/types'

type Props = {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSaved: () => void
}

type Item = { product_id: string; quantity: string; per_item_price: string }

const emptyItem: Item = { product_id: '', quantity: '1', per_item_price: '' }

export function OrderFormDialog({ open, onOpenChange, onSaved }: Props) {
  const [orderType, setOrderType] = useState('sale')
  const [partnerId, setPartnerId] = useState('')
  const [currency, setCurrency] = useState('1')
  const [exchangeRate, setExchangeRate] = useState('1')
  const [items, setItems] = useState<Item[]>([{ ...emptyItem }])
  const [saving, setSaving] = useState(false)

  const [products, setProducts] = useState<Product[]>([])
  const [partners, setPartners] = useState<Partner[] | null>(null)

  useEffect(() => {
    if (!open) return
    setOrderType('sale')
    setPartnerId('')
    setCurrency('1')
    setExchangeRate('1')
    setItems([{ ...emptyItem }])

    apiFetch<{ products: Product[] }>('/products/')
      .then((res) => setProducts(res.products ?? []))
      .catch(() => setProducts([]))

    apiFetch<{ partners: Partner[] }>('/partners/')
      .then((res) => setPartners(res.partners ?? []))
      .catch(() => setPartners(null))
  }, [open])

  function updateItem(index: number, key: keyof Item, value: string) {
    setItems((prev) =>
      prev.map((it, i) => (i === index ? { ...it, [key]: value } : it)),
    )
  }

  function addItem() {
    setItems((prev) => [...prev, { ...emptyItem }])
  }

  function removeItem(index: number) {
    setItems((prev) => prev.filter((_, i) => i !== index))
  }

  const total = items.reduce(
    (sum, it) => sum + Number(it.quantity || 0) * Number(it.per_item_price || 0),
    0,
  )

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!partnerId) {
      toast.error('شریک تجاری را انتخاب کنید.')
      return
    }
    if (items.some((it) => !it.product_id || !it.quantity || !it.per_item_price)) {
      toast.error('اطلاعات تمام اقلام سفارش را کامل کنید.')
      return
    }

    setSaving(true)
    try {
      await apiFetch('/orders/', {
        method: 'POST',
        body: {
          order_type: orderType,
          business_partner_name: Number(partnerId),
          currency: Number(currency),
          exchange_rate: Number(exchangeRate),
          order_items: items.map((it) => ({
            product_id: Number(it.product_id),
            quantity: Number(it.quantity),
            per_item_price: Number(it.per_item_price),
          })),
        },
      })
      toast.success('سفارش جدید ثبت شد.')
      onSaved()
      onOpenChange(false)
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در ثبت سفارش.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-2xl">
        <DialogHeader>
          <DialogTitle>ثبت سفارش جدید</DialogTitle>
          <DialogDescription>
            نوع سفارش، شریک تجاری و اقلام را مشخص کنید.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="flex flex-col gap-2">
              <Label htmlFor="order_type">نوع سفارش</Label>
              <Select value={orderType} onValueChange={setOrderType}>
                <SelectTrigger id="order_type">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {Object.entries(ORDER_TYPES).map(([value, label]) => (
                    <SelectItem key={value} value={value}>
                      {label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="partner">شریک تجاری</Label>
              {partners ? (
                <Select value={partnerId} onValueChange={setPartnerId}>
                  <SelectTrigger id="partner">
                    <SelectValue placeholder="انتخاب شریک" />
                  </SelectTrigger>
                  <SelectContent>
                    {partners.map((p) => (
                      <SelectItem key={p.id} value={String(p.id)}>
                        {p.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : (
                <Input
                  id="partner"
                  type="number"
                  min={1}
                  dir="ltr"
                  className="text-left"
                  value={partnerId}
                  onChange={(e) => setPartnerId(e.target.value)}
                  placeholder="شناسه شریک تجاری"
                />
              )}
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="currency">واحد پول</Label>
              <Select value={currency} onValueChange={setCurrency}>
                <SelectTrigger id="currency">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {CURRENCY_OPTIONS.map((c) => (
                    <SelectItem key={c.value} value={String(c.value)}>
                      {c.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor="exchange_rate">نرخ تبدیل</Label>
              <Input
                id="exchange_rate"
                type="number"
                min={0}
                step="0.01"
                dir="ltr"
                className="text-left"
                value={exchangeRate}
                onChange={(e) => setExchangeRate(e.target.value)}
              />
            </div>
          </div>

          <div className="flex flex-col gap-3 rounded-lg border border-border p-3">
            <div className="flex items-center justify-between">
              <Label>اقلام سفارش</Label>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={addItem}
              >
                <Plus className="size-4" />
                افزودن قلم
              </Button>
            </div>

            {items.map((item, index) => (
              <div
                key={index}
                className="grid grid-cols-1 items-end gap-2 sm:grid-cols-[1fr_auto_auto_auto]"
              >
                <div className="flex flex-col gap-1">
                  <span className="text-xs text-muted-foreground">محصول</span>
                  <Select
                    value={item.product_id}
                    onValueChange={(v) => updateItem(index, 'product_id', v)}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="انتخاب محصول" />
                    </SelectTrigger>
                    <SelectContent>
                      {products
                        .filter((p) => !p.is_frozen)
                        .map((p) => (
                          <SelectItem key={p.id} value={String(p.id)}>
                            {p.name}
                          </SelectItem>
                        ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex flex-col gap-1">
                  <span className="text-xs text-muted-foreground">تعداد</span>
                  <Input
                    type="number"
                    min={1}
                    dir="ltr"
                    className="w-20 text-left"
                    value={item.quantity}
                    onChange={(e) =>
                      updateItem(index, 'quantity', e.target.value)
                    }
                  />
                </div>

                <div className="flex flex-col gap-1">
                  <span className="text-xs text-muted-foreground">قیمت واحد</span>
                  <Input
                    type="number"
                    min={0}
                    dir="ltr"
                    className="w-32 text-left"
                    value={item.per_item_price}
                    onChange={(e) =>
                      updateItem(index, 'per_item_price', e.target.value)
                    }
                  />
                </div>

                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="text-destructive hover:text-destructive"
                  onClick={() => removeItem(index)}
                  disabled={items.length === 1}
                  aria-label="حذف قلم"
                >
                  <Trash2 className="size-4" />
                </Button>
              </div>
            ))}

            <div className="flex items-center justify-between border-t border-border pt-2 text-sm">
              <span className="text-muted-foreground">جمع کل</span>
              <span className="font-bold text-foreground">
                {formatPrice(total)}
              </span>
            </div>
          </div>

          <DialogFooter className="gap-2 sm:gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
            >
              انصراف
            </Button>
            <Button type="submit" disabled={saving}>
              {saving && <Loader2 className="size-4 animate-spin" />}
              ثبت سفارش
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
