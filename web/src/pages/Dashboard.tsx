function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

function MetricCard({
  label,
  value,
  icon,
  iconBg,
  iconColor,
  badge,
  trend,
  errorBorder,
  errorText,
}: {
  label: string;
  value: string;
  icon: string;
  iconBg: string;
  iconColor: string;
  badge?: { icon: string; label: string; value: string };
  trend?: boolean;
  errorBorder?: boolean;
  errorText?: boolean;
}) {
  const colors = { iconBg, iconColor };

  return (
    <div
      className={`group relative flex h-32 flex-col justify-between overflow-hidden rounded-xl border border-outline-variant bg-surface-container-lowest p-6 shadow-[0_4px_12px_rgba(0,0,0,0.02)] transition-all hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)] ${
        errorBorder ? "border-l-4 border-l-error" : ""
      }`}
    >
      <div className="flex items-start justify-between">
        <span className="font-label-caps text-label-caps tracking-wider text-on-surface-variant">
          {label}
        </span>
        <div className={`rounded p-2 ${colors.iconBg} ${colors.iconColor}`}>
          <Icon name={icon} />
        </div>
      </div>
      <div className="flex items-baseline space-x-2">
        <span
          className={`font-display-lg text-display-lg ${
            errorText ? "text-error" : "text-on-surface"
          }`}
        >
          {value}
        </span>
        {badge && (
          <span
            className={`flex items-center rounded px-1.5 py-0.5 text-xs font-bold ${
              trend
                ? "bg-secondary-container/20 text-secondary"
                : "bg-error-container/20 text-error"
            }`}
          >
            <Icon name={badge.icon} className="mr-0.5 text-[14px]" />
            {badge.value}
          </span>
        )}
        {!badge && errorText && (
          <span className="text-xs font-medium text-on-surface-variant">
            Requires Action
          </span>
        )}
      </div>
    </div>
  );
}

function TopSeller({
  initials,
  initialsBg,
  initialsColor,
  name,
  description,
  amount,
}: {
  initials: string;
  initialsBg: string;
  initialsColor: string;
  name: string;
  description: string;
  amount: string;
}) {
  return (
    <div className="flex items-center justify-between rounded-lg p-3 transition-colors hover:bg-surface-container-low">
      <div className="flex items-center">
        <div
          className={`mr-3 flex h-10 w-10 items-center justify-center rounded text-xs font-bold ${initialsBg} ${initialsColor}`}
        >
          {initials}
        </div>
        <div>
          <p className="text-sm font-bold">{name}</p>
          <p className="text-xs text-on-surface-variant">{description}</p>
        </div>
      </div>
      <span className="font-data-mono text-data-mono text-secondary">
        {amount}
      </span>
    </div>
  );
}

function TransactionRow({
  id,
  initials,
  name,
  amount,
  status,
  date,
}: {
  id: string;
  initials: string;
  name: string;
  amount: string;
  status: "Completed" | "Pending" | "Refunded";
  date: string;
}) {
  const statusStyles: Record<string, string> = {
    Completed: "bg-secondary-container/30 text-on-secondary-fixed-variant",
    Pending: "bg-tertiary-container/20 text-on-tertiary-fixed-variant",
    Refunded: "bg-error-container/40 text-on-error-container",
  };

  return (
    <tr className="transition-colors hover:bg-[#F8FAFF] group">
      <td className="font-data-mono px-6 py-4 text-sm">{id}</td>
      <td className="px-6 py-4">
        <div className="flex items-center">
          <div className="mr-3 flex h-7 w-7 items-center justify-center rounded-full bg-outline-variant/20 text-[10px] font-bold">
            {initials}
          </div>
          <span className="font-medium">{name}</span>
        </div>
      </td>
      <td className="px-6 py-4 text-right font-data-mono">{amount}</td>
      <td className="px-6 py-4">
        <div className="flex justify-center">
          <span
            className={`rounded-full px-2 py-1 text-[10px] font-bold uppercase tracking-wider ${statusStyles[status]}`}
          >
            {status}
          </span>
        </div>
      </td>
      <td className="px-6 py-4 text-xs text-on-surface-variant">{date}</td>
      <td className="px-6 py-4 text-right">
        <button className="rounded p-1 opacity-0 transition-colors hover:bg-surface-container-high group-hover:opacity-100">
          <Icon name="more_vert" className="text-lg" />
        </button>
      </td>
    </tr>
  );
}

export default function Dashboard() {
  return (
    <>
      {/* Header Section */}
      <div className="mb-8 flex items-end justify-between">
        <div>
          <h2 className="mb-0 font-display-lg text-display-lg text-on-surface">
            Sales Overview
          </h2>
          <p className="text-body-lg text-on-surface-variant">
            Real-time performance and inventory metrics.
          </p>
        </div>
        <div className="flex space-x-3">
          <button className="flex items-center rounded-lg border border-primary px-4 py-2 font-semibold text-primary transition-all hover:bg-surface-container-high">
            <Icon name="calendar_today" className="mr-2 text-sm" />
            Last 7 Days
          </button>
          <button className="btn-sky-action">
            <Icon name="download" className="mr-2 text-sm" />
            Export Report
          </button>
        </div>
      </div>

      {/* Metric Cards */}
      <div className="mb-8 grid grid-cols-1 gap-card-gap md:grid-cols-2 lg:grid-cols-4">
        <MetricCard
          label="DAILY REVENUE"
          value="$4,250"
          icon="payments"
          iconBg="bg-primary/5"
          iconColor="text-primary"
          badge={{ icon: "trending_up", label: "+12%", value: "12%" }}
          trend
        />
        <MetricCard
          label="ORDERS TODAY"
          value="84"
          icon="shopping_cart"
          iconBg="bg-secondary/5"
          iconColor="text-secondary"
          badge={{ icon: "trending_up", label: "+8%", value: "8%" }}
          trend
        />
        <MetricCard
          label="NEW CUSTOMERS"
          value="12"
          icon="person_add"
          iconBg="bg-tertiary/5"
          iconColor="text-tertiary"
          badge={{ icon: "trending_down", label: "-3%", value: "3%" }}
        />
        <MetricCard
          label="LOW STOCK ITEMS"
          value="05"
          icon="inventory"
          iconBg="bg-error/5"
          iconColor="text-error"
          errorBorder
          errorText
        />
      </div>

      {/* Middle Section */}
      <div className="mb-8 grid grid-cols-1 gap-card-gap lg:grid-cols-3">
        {/* Sales Chart */}
        <div className="rounded-xl border border-outline-variant bg-surface-container-lowest p-6 shadow-[0_4px_12px_rgba(0,0,0,0.02)] transition-all hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)] lg:col-span-2">
          <div className="mb-6 flex items-center justify-between">
            <h3 className="text-body-lg font-bold">Weekly Sales Trends</h3>
            <div className="flex items-center space-x-4">
              <div className="flex items-center">
                <span className="mr-2 h-3 w-3 rounded-full bg-primary" />
                <span className="text-xs text-on-surface-variant">Current Week</span>
              </div>
              <div className="flex items-center">
                <span className="mr-2 h-3 w-3 rounded-full bg-outline-variant" />
                <span className="text-xs text-on-surface-variant">Previous Week</span>
              </div>
            </div>
          </div>
          <div className="relative h-[300px] w-full">
            <div className="absolute inset-0 flex flex-col justify-between pb-8">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-0 w-full border-b border-surface-container" />
              ))}
            </div>
            <svg className="absolute inset-0 h-full w-full" preserveAspectRatio="none">
              <path
                d="M0 250 Q 150 200 300 150 T 600 100 T 900 180 T 1200 50"
                fill="none"
                stroke="#004bcd"
                strokeWidth="3"
              />
              <path
                d="M0 280 Q 150 240 300 210 T 600 180 T 900 230 T 1200 120"
                fill="none"
                stroke="#c3c5d9"
                strokeDasharray="4"
                strokeWidth="2"
              />
            </svg>
            <div className="absolute bottom-0 flex w-full justify-between px-2 text-[10px] font-medium uppercase tracking-tighter text-on-surface-variant">
              <span>Mon</span><span>Tue</span><span>Wed</span><span>Thu</span>
              <span>Fri</span><span>Sat</span><span>Sun</span>
            </div>
          </div>
        </div>

        {/* Top Sellers */}
        <div className="rounded-xl border border-outline-variant bg-surface-container-lowest p-6 shadow-[0_4px_12px_rgba(0,0,0,0.02)] transition-all hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)]">
          <div className="mb-6 flex items-center justify-between">
            <h3 className="text-body-lg font-bold">Top Sellers</h3>
            <button className="text-xs font-bold text-primary hover:underline">
              View All
            </button>
          </div>
          <div className="space-y-4">
            <TopSeller
              initials="AMX"
              initialsBg="bg-primary/5"
              initialsColor="text-primary"
              name="Amoxicillin 500mg"
              description="Antibiotic • 1.2k units"
              amount="+$3,420"
            />
            <TopSeller
              initials="PAR"
              initialsBg="bg-secondary/5"
              initialsColor="text-secondary"
              name="Paracetamol Extra"
              description="Analgesic • 840 units"
              amount="+$1,890"
            />
            <TopSeller
              initials="VIT"
              initialsBg="bg-tertiary/5"
              initialsColor="text-tertiary"
              name="Vitamin C 1000mg"
              description="Supplement • 620 units"
              amount="+$1,140"
            />
            <div className="mt-2 border-t border-surface-container pt-4">
              <p className="text-xs leading-relaxed italic text-on-surface-variant">
                "Inventory levels for top sellers are within safety thresholds
                for the next 48 hours."
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Transactions */}
      <div className="overflow-hidden rounded-xl border border-outline-variant bg-surface-container-lowest shadow-[0_4px_12px_rgba(0,0,0,0.02)] transition-all hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)]">
        <div className="flex items-center justify-between border-b border-outline-variant bg-surface-container-low/30 px-6 py-4">
          <h3 className="text-body-lg font-bold">Recent Transactions</h3>
          <div className="relative">
            <Icon
              name="filter_list"
              className="absolute left-2 top-1/2 -translate-y-1/2 text-xs text-on-surface-variant"
            />
            <input
              className="rounded-lg border border-outline-variant bg-surface-container-low py-1.5 pl-8 pr-3 text-xs outline-none focus:border-primary focus:ring-2 focus:ring-primary"
              placeholder="Filter transactions..."
              type="text"
            />
          </div>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="bg-surface-container-low/50">
                <th className="px-6 py-4 font-label-caps text-label-caps uppercase tracking-wider text-on-surface-variant">
                  Order ID
                </th>
                <th className="px-6 py-4 font-label-caps text-label-caps uppercase tracking-wider text-on-surface-variant">
                  Customer
                </th>
                <th className="px-6 py-4 text-right font-label-caps text-label-caps uppercase tracking-wider text-on-surface-variant">
                  Amount
                </th>
                <th className="px-6 py-4 text-center font-label-caps text-label-caps uppercase tracking-wider text-on-surface-variant">
                  Status
                </th>
                <th className="px-6 py-4 font-label-caps text-label-caps uppercase tracking-wider text-on-surface-variant">
                  Date
                </th>
                <th className="px-6 py-4" />
              </tr>
            </thead>
            <tbody className="divide-y divide-outline-variant/30">
              <TransactionRow
                id="#PF-9842"
                initials="JD"
                name="John Doe"
                amount="$124.50"
                status="Completed"
                date="Oct 12, 14:32"
              />
              <TransactionRow
                id="#PF-9841"
                initials="EM"
                name="Elena Martinez"
                amount="$82.00"
                status="Pending"
                date="Oct 12, 14:15"
              />
              <TransactionRow
                id="#PF-9840"
                initials="RW"
                name="Robert Wilson"
                amount="$45.10"
                status="Completed"
                date="Oct 12, 13:58"
              />
              <TransactionRow
                id="#PF-9839"
                initials="AK"
                name="Alice Kim"
                amount="$312.00"
                status="Completed"
                date="Oct 12, 13:42"
              />
              <TransactionRow
                id="#PF-9838"
                initials="TH"
                name="Thomas Hardy"
                amount="$18.45"
                status="Refunded"
                date="Oct 12, 13:10"
              />
            </tbody>
          </table>
        </div>
      </div>
    </>
  );
}
