import type { Permission } from './permissions'

export type TransitionAction = {
  action: string
  label: string
  permission: Permission
  nextStatus: string
}

const SALE_FLOW: Record<string, TransitionAction[]> = {
  Pending: [
    { action: 'approve', label: 'تأیید سفارش', permission: 'orders.update', nextStatus: 'Approved' },
    { action: 'cancel', label: 'لغو سفارش', permission: 'orders.update', nextStatus: 'Canceled' },
  ],
  Approved: [
    { action: 'pack', label: 'بسته‌بندی', permission: 'orders.pack', nextStatus: 'Packed' },
    { action: 'cancel', label: 'لغو سفارش', permission: 'orders.update', nextStatus: 'Canceled' },
  ],
  Packed: [
    { action: 'ship', label: 'ارسال', permission: 'orders.ship', nextStatus: 'Shipped' },
    { action: 'cancel', label: 'لغو سفارش', permission: 'orders.update', nextStatus: 'Canceled' },
  ],
  Shipped: [],
  Canceled: [],
}

const PURCHASE_FLOW: Record<string, TransitionAction[]> = {
  Pending: [
    { action: 'approve', label: 'تأیید سفارش', permission: 'orders.update', nextStatus: 'Approved' },
    { action: 'cancel', label: 'لغو سفارش', permission: 'orders.update', nextStatus: 'Canceled' },
  ],
  Approved: [
    { action: 'wait', label: 'انتظار دریافت', permission: 'orders.update', nextStatus: 'Waiting' },
    { action: 'cancel', label: 'لغو سفارش', permission: 'orders.update', nextStatus: 'Canceled' },
  ],
  Waiting: [
    { action: 'receive', label: 'دریافت کالا', permission: 'orders.receive', nextStatus: 'Received' },
  ],
  Received: [],
  Canceled: [],
}

export function getTransitions(
  orderType: string,
  status: string,
): TransitionAction[] {
  const flow = orderType === 'purchase' ? PURCHASE_FLOW : SALE_FLOW
  return flow[status] ?? []
}
