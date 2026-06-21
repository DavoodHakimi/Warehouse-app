'use client'

import { useTheme } from 'next-themes'
import { Monitor, Moon, Sun } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { buttonVariants } from '@/components/ui/button'
import type { VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'
import { useEffect, useState } from 'react'

const THEMES = [
  { value: 'system', label: 'خودکار', Icon: Monitor },
  { value: 'light', label: 'روشن', Icon: Sun },
  { value: 'dark', label: 'تاریک', Icon: Moon },
] as const

export function ThemeToggle() {
  const { theme, setTheme } = useTheme()
  const [mounted, setMounted] = useState(false)

  useEffect(() => setMounted(true), [])

  const current = THEMES.find((t) => t.value === theme) ?? THEMES[0]
  const Icon = current.Icon

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        className={cn(
          buttonVariants({ variant: 'ghost', size: 'icon' }),
          'size-9'
        )}
        aria-label="تغییر تم"
      >
        {mounted ? (
          <Icon className="size-[1.2rem]" />
        ) : (
          <Sun className="size-[1.2rem]" />
        )}
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start">
        {THEMES.map(({ value, label, Icon }) => (
          <DropdownMenuItem
            key={value}
            onClick={() => setTheme(value)}
            className={theme === value ? 'font-bold' : ''}
          >
            <Icon className="size-4" />
            {label}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
