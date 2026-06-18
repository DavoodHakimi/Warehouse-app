export type Permission =
  | 'users.read'
  | 'users.create'
  | 'users.update'
  | 'users.delete'
  | 'partners.read'
  | 'partners.create'
  | 'partners.update'
  | 'partners.delete'
  | 'products.read'
  | 'products.create'
  | 'products.update'
  | 'products.delete'
  | 'orders.read'
  | 'orders.create'
  | 'orders.update'
  | 'orders.delete'
  | 'orders.receive'
  | 'orders.ship'
  | 'orders.pack'

const MATRIX: Record<Permission, number[]> = {
  'users.read': [1],
  'users.create': [1],
  'users.update': [1],
  'users.delete': [1],
  'partners.read': [1],
  'partners.create': [1],
  'partners.update': [1],
  'partners.delete': [1],
  'products.read': [1, 2, 3, 4, 5],
  'products.create': [1, 2],
  'products.update': [1, 2],
  'products.delete': [1, 2],
  'orders.read': [1, 2, 3, 4, 5],
  'orders.create': [1, 2],
  'orders.update': [1, 2],
  'orders.delete': [1, 2],
  'orders.receive': [1, 2, 3, 4],
  'orders.ship': [1, 2, 3, 5],
  'orders.pack': [1, 2, 3],
}

export function can(role: number | undefined, permission: Permission): boolean {
  if (!role) return false
  return MATRIX[permission]?.includes(role) ?? false
}
