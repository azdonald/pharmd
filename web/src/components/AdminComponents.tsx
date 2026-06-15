import type { ReactNode } from "react";

export function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

export function PageHeader({
  title,
  description,
  actions,
}: {
  title: string;
  description: string;
  actions?: ReactNode;
}) {
  return (
    <div className="mb-8 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
      <div>
        <h2 className="mb-0 font-display-lg text-display-lg text-on-surface">{title}</h2>
        <p className="text-body-lg text-on-surface-variant">{description}</p>
      </div>
      {actions ? <div className="flex flex-wrap gap-3 items-center">{actions}</div> : null}
    </div>
  );
}

export function Panel({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <div className={`overflow-hidden rounded-xl border border-outline-variant bg-surface-container-lowest shadow-[0_4px_12px_rgba(0,0,0,0.02)] transition-all hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)] ${className}`}>
      {children}
    </div>
  );
}
