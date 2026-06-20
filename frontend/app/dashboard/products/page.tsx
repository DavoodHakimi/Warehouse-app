'use client'

import { useState } from 'react'
import { Plus, Pencil, Trash2, Snowflake } from 'lucide-react'
import { toast } from 'sonner'
import { PageHeader } from '@/components/page-header'
import { DataError, EmptyState, TableSkeleton } from '@/components/data-states'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { ProductFormDialog } from './product-form-dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { useApi } from '@/lib/use-api'
import { apiFetch, ApiError } from '@/lib/api'
import { useAuth } from '@/lib/auth'
import { formatPrice } from '@/lib/utils'
import type { Product } from '@/lib/types'

export default function ProductsPage() {
  const { can } = useAuth()
  const { data, error, loading, refetch } =
    useApi<{ products: Product[] }>('/products/')
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Product | null>(null)
  const [deleting, setDeleting] = useState<Product | null>(null)

  const products = data?.products ?? []
  const canCreate = can('products.create')
  const canUpdate = can('products.update')
  const canDelete = can('products.delete')

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(product: Product) {
    setEditing(product)
    setFormOpen(true)
  }

  async function handleDelete() {
    if (!deleting) return
    try {
      await apiFetch(`/products/${deleting.product_number}/`, {
        method: 'DELETE',
      })
      toast.success('محصول حذف شد.')
      setDeleting(null)
      refetch()
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در حذف.')
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="محصولات"
        description="مدیریت کالاها و موجودی انبار."
        action={
          canCreate ? (
            <Button onClick={openCreate}>
              <Plus className="size-4" />
              افزودن محصول
            </Button>
          ) : null
        }
      />

      {loading ? (
        <TableSkeleton columns={4} />
      ) : error ? (
        <DataError message={error} onRetry={refetch} />
      ) : products.length === 0 ? (
        <EmptyState
          title="محصولی یافت نشد"
          description="هنوز محصولی ثبت نشده است."
        />
      ) : (
        <Card className="overflow-hidden p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="text-right">نام محصول</TableHead>
                <TableHead className="text-right">شماره محصول</TableHead>
                <TableHead className="text-right">قیمت پایه</TableHead>
                <TableHead className="text-right">وضعیت</TableHead>
                {(canUpdate || canDelete) && (
                  <TableHead className="text-right">عملیات</TableHead>
                )}
              </TableRow>
            </TableHeader>
            <TableBody>
              {products.map((product) => (
                <TableRow key={product.id}>
                  <TableCell className="font-medium">{product.name}</TableCell>
                  <TableCell dir="ltr" className="text-right font-mono text-sm">
                    {product.product_number}
                  </TableCell>
                  <TableCell>{formatPrice(product.default_price)}</TableCell>
                  <TableCell>
                    {product.is_frozen ? (
                      <Badge variant="warning" className="gap-1">
                        <Snowflake className="size-3" />
                        منجمد
                      </Badge>
                    ) : (
                      <Badge variant="success">فعال</Badge>
                    )}
                  </TableCell>
                  {(canUpdate || canDelete) && (
                    <TableCell className="text-left">
                      <div className="flex justify-start gap-1">
                        {canUpdate && (
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => openEdit(product)}
                            aria-label="ویرایش"
                          >
                            <Pencil className="size-4" />
                          </Button>
                        )}
                        {canDelete && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="text-destructive hover:text-destructive"
                            onClick={() => setDeleting(product)}
                            aria-label="حذف"
                          >
                            <Trash2 className="size-4" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}

      <ProductFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        product={editing}
        onSaved={refetch}
      />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(o) => !o && setDeleting(null)}
        title="حذف محصول"
        description={`آیا از حذف «${deleting?.name}» مطمئن هستید؟`}
        confirmText="حذف"
        destructive
        onConfirm={handleDelete}
      />
    </div>
  )
}
