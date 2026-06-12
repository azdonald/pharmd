import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listStock, listAlerts, listExpiring, type StockItem, type StockAlert, type ExpiringBatch } from "../api/inventory";
import { listLocations, type Location } from "../api/locations";

function Icon({ name, className }: { name: string; className?: string }) {
  return <span className={`material-symbols-outlined ${className ?? ""}`}>{name}</span>;
}

export default function Inventory() {
  const [locations, setLocations] = useState<Location[]>([]);
  const [locationId, setLocationId] = useState("");
  const [stock, setStock] = useState<StockItem[]>([]);
  const [alerts, setAlerts] = useState<StockAlert[]>([]);
  const [expiring, setExpiring] = useState<ExpiringBatch[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 15;

  useEffect(() => {
    listLocations().then(res => {
      setLocations(res.data);
      if (res.data.length > 0) setLocationId(res.data[0].id);
    }).catch(console.error);
  }, []);

  const load = () => {
    if (!locationId) return;
    setLoading(true);
    Promise.all([
      listStock(locationId, page, limit, search),
      listAlerts(locationId),
      listExpiring(locationId, 30),
    ]).then(([s, a, e]) => {
      setStock(s.data);
      setTotal(s.total);
      setAlerts(a);
      setExpiring(e);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(load, [locationId, page, search]);

  const totalPages = Math.ceil(total / limit);

  return (
    <div>
      {/* Header */}
      <div className="flex justify-between items-end mb-8">
        <div>
          <div className="flex items-center gap-4 mb-1">
            <h2 className="font-display-lg text-display-lg text-on-surface">
              Inventory Management
            </h2>
            <select
              value={locationId}
              onChange={e => { setLocationId(e.target.value); setPage(1); }}
              className="rounded-lg border border-outline-variant bg-surface-container-lowest px-3 py-1.5 text-sm focus:ring-2 focus:ring-primary outline-none"
            >
              {locations.map(l => (
                <option key={l.id} value={l.id}>{l.name}</option>
              ))}
            </select>
          </div>
          <p className="text-body-lg text-on-surface-variant">
            Real-time medication stock tracking and replenishment
          </p>
        </div>
        <div className="flex space-x-3">
          <button className="flex items-center px-4 py-2 border border-primary text-primary font-semibold rounded-lg hover:bg-surface-container-high transition-all">
            <Icon name="print" className="mr-2" />
            Print Barcodes
          </button>
          <Link
            to="/app/inventory/receive"
            className="flex items-center px-4 py-2 bg-primary text-on-primary font-semibold rounded-lg hover:bg-primary-container shadow-md transition-all"
          >
            <Icon name="add" className="mr-2" />
            Receive Stock
          </Link>
        </div>
      </div>

      {/* Search */}
      <div className="mb-8 max-w-md">
        <div className="relative">
          <Icon name="search" className="absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant" />
          <input
            value={search}
            onChange={e => { setSearch(e.target.value); setPage(1); }}
            className="w-full rounded-lg border border-outline-variant bg-surface-container-lowest py-2 pl-10 pr-4 text-sm outline-none focus:ring-2 focus:ring-primary"
            placeholder="Search medications, categories..."
            type="text"
          />
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-card-gap mb-8">
        <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant group hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)] transition-all">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-primary-container/10 rounded-lg text-primary">
              <Icon name="inventory" />
            </div>
            {total > 0 && (
              <span className="bg-secondary-container text-on-secondary-container text-label-caps px-2 py-1 rounded-full">
                Total
              </span>
            )}
          </div>
          <p className="text-label-caps text-on-surface-variant uppercase mb-1">Total SKUs</p>
          <h3 className="font-headline-md text-headline-md text-on-surface">
            {loading ? "..." : total.toLocaleString()}
          </h3>
        </div>
        <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)] transition-all">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-error-container/20 rounded-lg text-error">
              <Icon name="warning" />
            </div>
            <span className="bg-error-container text-on-error-container text-label-caps px-2 py-1 rounded-full">
              {alerts.length > 0 ? "Critical" : "OK"}
            </span>
          </div>
          <p className="text-label-caps text-on-surface-variant uppercase mb-1">Low Stock Alerts</p>
          <h3 className="font-headline-md text-headline-md text-on-surface">
            {loading ? "..." : `${alerts.length} Items`}
          </h3>
        </div>
        <div className="bg-surface-container-lowest p-6 rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant hover:shadow-[0_12px_24px_rgba(13,97,255,0.06)] transition-all">
          <div className="flex justify-between items-start mb-4">
            <div className="p-3 bg-tertiary-fixed rounded-lg text-tertiary">
              <Icon name="event_busy" />
            </div>
            <span className="bg-tertiary-container text-on-tertiary-container text-label-caps px-2 py-1 rounded-full">
              30 Days
            </span>
          </div>
          <p className="text-label-caps text-on-surface-variant uppercase mb-1">Expiring Soon</p>
          <h3 className="font-headline-md text-headline-md text-on-surface">
            {loading ? "..." : `${expiring.length} Batches`}
          </h3>
        </div>
      </div>

      {/* Inventory Table */}
      <div className="bg-surface-container-lowest rounded-xl shadow-[0_4px_12px_rgba(0,0,0,0.02)] border border-outline-variant overflow-hidden">
        <div className="px-6 py-4 border-b border-outline-variant flex justify-between items-center bg-surface-container-low/30">
          <div className="flex space-x-4">
            <button className="text-primary font-semibold text-body-md border-b-2 border-primary pb-4 -mb-4">
              All Products
            </button>
          </div>
          <button className="flex items-center text-on-surface-variant hover:text-primary transition-all">
            <Icon name="filter_list" className="mr-1" />
            <span className="text-body-md">Filters</span>
          </button>
        </div>

        <div className="overflow-x-auto">
          {loading ? (
            <div className="p-12 text-center">
              <p className="text-on-surface-variant">Loading inventory...</p>
            </div>
          ) : stock.length === 0 ? (
            <div className="p-12 text-center">
              <p className="text-on-surface-variant">No products found</p>
            </div>
          ) : (
            <table className="w-full text-left">
              <thead>
                <tr className="bg-surface-container-low/50">
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Product</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Category</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Stock Level</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider text-center">Status</th>
                  <th className="px-6 py-4 font-label-caps text-label-caps text-on-surface-variant uppercase tracking-wider">Batches</th>
                  <th className="px-6 py-4"></th>
                </tr>
              </thead>
              <tbody className="divide-y divide-outline-variant/30">
                {stock.map(item => {
                  const isLow = item.total_quantity < item.reorder_level;
                  const isOut = item.total_quantity === 0;

                  let statusLabel = "In Stock";
                  let statusClass = "bg-secondary-container/20 text-secondary";
                  let barColor = "bg-primary";
                  if (isOut) {
                    statusLabel = "Out of Stock";
                    statusClass = "bg-outline-variant/20 text-on-surface-variant";
                    barColor = "bg-outline-variant";
                  } else if (isLow) {
                    statusLabel = "Low Stock";
                    statusClass = "bg-error-container/20 text-error";
                    barColor = "bg-error";
                  }

                  const barWidth = Math.min(
                    item.reorder_level > 0
                      ? Math.round((item.total_quantity / (item.reorder_level * 2)) * 100)
                      : 100,
                    100
                  );

                  return (
                    <tr key={item.product_id} className="hover:bg-surface-container-high/20 transition-colors group">
                      <td className="px-6 py-4">
                        <div className="flex flex-col">
                          <span className="font-semibold text-on-surface">
                            <Link to={`/app/products/${item.product_id}`} className="text-on-surface hover:text-primary transition-colors">
                              {item.product_name}
                            </Link>
                          </span>
                          <span className="text-label-caps text-on-surface-variant">
                            {item.brand_name || item.generic_name || item.classification || ""}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <span className="text-body-md text-on-surface-variant">{item.classification || "—"}</span>
                      </td>
                      <td className="px-6 py-4 w-56">
                        <div className="flex items-center space-x-2">
                          <div className="flex-1 h-2 bg-surface-container-high rounded-full overflow-hidden">
                            <div className={`h-full ${barColor} rounded-full transition-all`} style={{ width: `${barWidth}%` }} />
                          </div>
                          <span className="font-data-mono text-data-mono text-on-surface shrink-0">
                            {item.total_quantity}/{item.reorder_level}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-center">
                        <span className={`inline-flex items-center px-3 py-1 rounded-full text-label-caps font-bold ${statusClass}`}>
                          {statusLabel}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <details className="group">
                          <summary className="cursor-pointer text-body-md text-on-surface-variant hover:text-primary list-none flex items-center">
                            <Icon name="chevron_right" className="transition-transform group-open:rotate-90 text-sm" />
                            <span className="ml-1">{item.batches.length} batch(es)</span>
                          </summary>
                          <div className="mt-2 bg-surface-container-low rounded-lg p-3 space-y-2">
                            {item.batches.map(b => (
                              <div key={b.id} className="flex justify-between text-xs">
                                <span className="font-data-mono">{b.batch_number || "—"}</span>
                                <span>{b.quantity} units</span>
                                <span>${b.selling_price?.toFixed(2)}</span>
                                <span className={b.expiry_date && new Date(b.expiry_date) < new Date(Date.now() + 30 * 86400000) ? "text-error" : ""}>
                                  {b.expiry_date || "—"}
                                </span>
                              </div>
                            ))}
                          </div>
                        </details>
                      </td>
                      <td className="px-6 py-4 text-right">
                        <Link
                          to={`/app/products/${item.product_id}`}
                          className="text-on-surface-variant hover:text-primary"
                        >
                          <Icon name="more_vert" />
                        </Link>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          )}
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t border-outline-variant flex justify-between items-center bg-surface-container-low/10">
            <p className="text-body-md text-on-surface-variant">
              Showing {((page - 1) * limit) + 1} to {Math.min(page * limit, total)} of {total} items
            </p>
            <div className="flex space-x-2">
              <button
                onClick={() => setPage(p => Math.max(1, p - 1))}
                disabled={page === 1}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Icon name="chevron_left" />
              </button>
              {Array.from({ length: Math.min(totalPages, 5) }, (_, i) => {
                const start = Math.max(1, Math.min(page - 2, totalPages - 4));
                const p = start + i;
                if (p > totalPages) return null;
                return (
                  <button
                    key={p}
                    onClick={() => setPage(p)}
                    className={`w-10 h-10 rounded-md flex items-center justify-center font-bold text-sm ${
                      p === page
                        ? "bg-primary text-on-primary"
                        : "border border-outline-variant hover:bg-surface-container-high"
                    }`}
                  >
                    {p}
                  </button>
                );
              })}
              <button
                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="p-2 border border-outline-variant rounded-md hover:bg-surface-container-high disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <Icon name="chevron_right" />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
