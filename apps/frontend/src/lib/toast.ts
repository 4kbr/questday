import { toast } from 'sonner'
import { ApiError } from '@/apis/client'

export function toastSuccess(message: string) {
  toast.success(message)
}

/** Toast error dari hasil mutation. Pesan dari backend (ApiError.message);
 *  error jaringan → pesan generik. Tak pernah mengarang detail. */
export function toastApiError(err: unknown, fallback = 'Terjadi kesalahan') {
  if (err instanceof ApiError) {
    // status 0 = tak ada response (jaringan / server mati)
    toast.error(
      err.status === 0 ? 'Tidak dapat terhubung ke server' : err.message,
    )
    return
  }
  toast.error(fallback)
}
