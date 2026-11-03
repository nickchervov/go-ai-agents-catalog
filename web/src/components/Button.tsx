import type { ButtonHTMLAttributes } from 'react';

type Variant = 'default' | 'primary' | 'danger' | 'ghost';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: 'md' | 'sm';
  block?: boolean;
}

const VARIANT_CLASS: Record<Variant, string> = {
  default: '',
  primary: 'btn--primary',
  danger: 'btn--danger',
  ghost: 'btn--ghost',
};

export function Button({
  variant = 'default',
  size = 'md',
  block = false,
  className = '',
  ...rest
}: ButtonProps) {
  const classes = [
    'btn',
    VARIANT_CLASS[variant],
    size === 'sm' ? 'btn--sm' : '',
    block ? 'btn--block' : '',
    className,
  ]
    .filter(Boolean)
    .join(' ');
  return <button className={classes} {...rest} />;
}

export function FieldError({ children }: { children?: string }) {
  if (!children) return null;
  return <span className="field__error">{children}</span>;
}
