import { Component, type ErrorInfo, type ReactNode } from 'react'
import { Button } from '@/components/ui/button'

type Props = { children: ReactNode }
type State = { error: Error | null }

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null }

  static getDerivedStateFromError(error: Error): State {
    return { error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    if (import.meta.env.DEV) console.error('ErrorBoundary', error, info)
  }

  render() {
    if (!this.state.error) return this.props.children
    return (
      <div className="flex min-h-svh flex-col items-center justify-center gap-4 p-6 text-center">
        <h1 className="text-lg font-semibold">Ada yang salah</h1>
        <p className="text-sm text-muted-foreground">
          Coba muat ulang halaman.
        </p>
        {import.meta.env.DEV && (
          <pre className="max-w-lg overflow-auto rounded bg-muted p-3 text-left text-xs">
            {this.state.error.message}
          </pre>
        )}
        <Button onClick={() => window.location.reload()}>Muat ulang</Button>
      </div>
    )
  }
}
