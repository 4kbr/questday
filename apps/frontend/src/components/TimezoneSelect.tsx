import { useMemo, useState } from 'react'
import { Check, ChevronsUpDown } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList,
} from '@/components/ui/command'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { COMMON_TIMEZONES } from '@/lib/timezones'

type Props = {
  value: string
  onChange: (tz: string) => void
  id?: string
  'aria-label'?: string
}

/**
 * Combobox timezone IANA. Ditaruh di `components/` (bukan `features/auth/`)
 * karena dipakai `RegisterForm` DAN `ProfileForm`. Daftar penuh dari
 * `Intl.supportedValuesOf('timeZone')` (~400 zona); fallback daftar pendek
 * untuk browser lama.
 */
export function TimezoneSelect({
  value,
  onChange,
  id,
  'aria-label': ariaLabel,
}: Props) {
  const [open, setOpen] = useState(false)

  const zones = useMemo<string[]>(() => {
    const supported = (
      Intl as unknown as {
        supportedValuesOf?: (key: string) => string[]
      }
    ).supportedValuesOf
    return typeof supported === 'function'
      ? supported('timeZone')
      : [...COMMON_TIMEZONES]
  }, [])

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          id={id}
          type="button"
          variant="outline"
          role="combobox"
          aria-expanded={open}
          aria-label={ariaLabel}
          className="w-full justify-between font-normal"
        >
          {value || 'Pilih timezone'}
          <ChevronsUpDown className="ml-2 size-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[--radix-popover-trigger-width] p-0"
        align="start"
      >
        <Command>
          <CommandInput placeholder="Cari timezone..." />
          <CommandList className="max-h-72 overflow-y-auto">
            <CommandEmpty>Timezone tidak ditemukan.</CommandEmpty>
            {zones.map((tz) => (
              <CommandItem
                key={tz}
                value={tz}
                onSelect={() => {
                  onChange(tz)
                  setOpen(false)
                }}
              >
                <Check
                  className={cn(
                    'mr-2 size-4',
                    tz === value ? 'opacity-100' : 'opacity-0',
                  )}
                />
                {tz}
              </CommandItem>
            ))}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
