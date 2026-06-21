export type User = {
  id: number
  full_name: string
  user_name: string
  user_type_id: number
  phone_number: string
  email: string
}

export type Partner = {
  id: number
  name: string
  business_partner_type: string
  phone_number: string
  email: string
  contact_name: string
  contact_phone_number: string
}

export type Product = {
  id: number
  name: string
  product_number: string
  is_frozen: boolean
  default_price: number
}
export type OrderItem = {
  product_id: number
  quantity: number
  per_item_price: number
}

export type Order = {
  id: number
  order_type: string
  order_number: string
  status: string
  business_partner_name: string
  currency: string
  exchange_rate: number
  order_items?: OrderItem[]
}
