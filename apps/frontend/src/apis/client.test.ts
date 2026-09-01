import AxiosMockAdapter from 'axios-mock-adapter'
import { afterEach, beforeEach, describe, expect, it } from 'vitest'
import { ApiError, api, setTokenGetter } from '@/apis/client'

let mock: AxiosMockAdapter

beforeEach(() => {
  mock = new AxiosMockAdapter(api)
})

afterEach(() => {
  mock.restore()
  // Kembalikan token getter ke keadaan "tidak ada token".
  setTokenGetter(() => null)
})

describe('apis/client', () => {
  it('maps a 500 enveloped error to ApiError (status + code)', async () => {
    mock.onGet('/x').reply(500, {
      error: { code: 'internal_error', message: 'boom' },
    })

    await expect(api.get('/x')).rejects.toMatchObject({
      name: 'ApiError',
      status: 500,
      code: 'internal_error',
      message: 'boom',
    })
    await expect(api.get('/x')).rejects.toBeInstanceOf(ApiError)
  })

  it('maps a 409 enveloped error to ApiError with the domain code', async () => {
    mock.onGet('/x').reply(409, {
      error: { code: 'already_completed', message: 'quest sudah selesai' },
    })

    await expect(api.get('/x')).rejects.toMatchObject({
      status: 409,
      code: 'already_completed',
    })
  })

  it('returns raw success bodies without unwrapping a {data} envelope', async () => {
    mock.onGet('/x').reply(200, { id: '1', title: 'x' })

    const res = await api.get('/x')
    expect(res.data).toEqual({ id: '1', title: 'x' })
  })

  it('attaches Authorization: Bearer <token> from the injected token getter', async () => {
    setTokenGetter(() => 'tok')
    let seen: string | undefined
    mock.onGet('/x').reply((config) => {
      seen = config.headers?.Authorization as string | undefined
      return [200, {}]
    })

    await api.get('/x')
    expect(seen).toBe('Bearer tok')
  })
})
