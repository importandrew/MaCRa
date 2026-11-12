import { Component, type ReactNode } from 'react';
import { reportError } from '../../monitoring';
import CrashState from './CrashState';

type Props = {
  // resetKey clears a caught error when it changes: RootLayout passes the
  // path, so going to another page recovers.
  resetKey: string;
  children: ReactNode;
};

type State = {
  error: Error | null;
  // errorKey is the resetKey the error was caught under.
  errorKey: string | null;
};

// ErrorBoundary catches a render error in the page below it and shows
// CrashState instead of unmounting the whole app (a blank white screen).
// Errors in event handlers and async code don't reach it; those show toasts
// (Request.ts). React has no hook for this, so it is a class.
export default class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null, errorKey: null };

  static getDerivedStateFromError(error: Error): Partial<State> {
    return { error };
  }

  static getDerivedStateFromProps(props: Props, state: State): Partial<State> {
    if (state.error && state.errorKey === null) {
      return { errorKey: props.resetKey };
    }
    if (state.error && state.errorKey !== props.resetKey) {
      return { error: null, errorKey: null };
    }
    return {};
  }

  // componentDidCatch reports the error: Hawk's global handlers don't see
  // errors a boundary catches. In development React logs it with its
  // component stack itself.
  componentDidCatch(error: Error) {
    reportError(error);
  }

  render() {
    if (this.state.error) {
      return <CrashState error={this.state.error} />;
    }
    return this.props.children;
  }
}
