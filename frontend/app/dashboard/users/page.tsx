'use client'

import { useState } from 'react'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { PageHeader } from '@/components/page-header'
import { DataError, EmptyState, TableSkeleton } from '@/components/data-states'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { UserFormDialog } from './user-form-dialog'
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
import { ROLES } from '@/lib/constants'
import type { User } from '@/lib/types'

export default function UsersPage() {
  const { can } = useAuth()
  const { data, error, loading, refetch } = useApi<{ users: User[] }>('/users/')
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)
  const [deleting, setDeleting] = useState<User | null>(null)

  const users = data?.users ?? []
  const canCreate = can('users.create')
  const canUpdate = can('users.update')
  const canDelete = can('users.delete')

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(user: User) {
    setEditing(user)
    setFormOpen(true)
  }

  async function handleDelete() {
    if (!deleting) return
    try {
      await apiFetch(`/users/${deleting.id}`, { method: 'DELETE' })
      toast.success('کاربر حذف شد.')
      setDeleting(null)
      refetch()
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در حذف کاربر.')
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="مدیریت کاربران"
        description="افزودن، ویرایش و حذف کاربران شرکت."
        action={
          canCreate ? (
            <Button onClick={openCreate}>
              <Plus className="size-4" />
              افزودن کاربر
            </Button>
          ) : null
        }
      />

      {loading ? (
        <TableSkeleton columns={5} />
      ) : error ? (
        <DataError message={error} onRetry={refetch} />
      ) : users.length === 0 ? (
        <EmptyState
          title="کاربری یافت نشد"
          description="هنوز کاربری ثبت نشده است."
        />
      ) : (
        <Card className="overflow-hidden p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="text-right">نام و نام خانوادگی</TableHead>
                <TableHead className="text-right">نام کاربری</TableHead>
                <TableHead className="text-right">نقش</TableHead>
                <TableHead className="text-right">شماره تماس</TableHead>
                <TableHead className="text-right">ایمیل</TableHead>
                {(canUpdate || canDelete) && (
                  <TableHead className="text-left">عملیات</TableHead>
                )}
              </TableRow>
            </TableHeader>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell className="font-medium">{user.full_name}</TableCell>
                  <TableCell dir="ltr" className="text-right">
                    {user.user_name}
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary">
                      {ROLES[user.user_type_id] ?? '—'}
                    </Badge>
                  </TableCell>
                  <TableCell dir="ltr" className="text-right">
                    {user.phone_number}
                  </TableCell>
                  <TableCell dir="ltr" className="text-right">
                    {user.email}
                  </TableCell>
                  {(canUpdate || canDelete) && (
                    <TableCell className="text-left">
                      <div className="flex justify-start gap-1">
                        {canUpdate && (
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => openEdit(user)}
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
                            onClick={() => setDeleting(user)}
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

      <UserFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        user={editing}
        onSaved={refetch}
      />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(o) => !o && setDeleting(null)}
        title="حذف کاربر"
        description={`آیا از حذف «${deleting?.full_name}» مطمئن هستید؟ این عملیات قابل بازگشت نیست.`}
        confirmText="حذف"
        destructive
        onConfirm={handleDelete}
      />
    </div>
  )
}
