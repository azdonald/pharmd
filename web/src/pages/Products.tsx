import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { listProducts, listCategories, deleteProduct, type Product, type ProductCategory } from "../api/products";

export default function Products() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<ProductCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState(searchParams.get("query") || "");

  const page = parseInt(searchParams.get("page") || "1");
  const query = searchParams.get("query") || "";
  const categoryId = searchParams.get("category_id") || "";

  const load = () => {
    setLoading(true);
    Promise.all([
      listProducts(page, 20, query, categoryId),
      listCategories(),
    ]).then(([res]) => {
      setProducts(res.data);
      setTotal(res.total);
    }).catch(console.error).finally(() => setLoading(false));
  };

  useEffect(() => {
    listCategories().then(setCategories).catch(console.error);
  }, []);

  useEffect(load, [page, query, categoryId]);

  const catMap = Object.fromEntries(categories.map(c => [c.id, c.name]));

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearchParams({ query: search, page: "1" });
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Delete this product?")) return;
    try {
      await deleteProduct(id);
      load();
    } catch (err) {
      alert(err instanceof Error ? err.message : "Delete failed");
    }
  };

  const totalPages = Math.ceil(total / 20);

  return (
    <div>
      <div className="page-header">
        <h1>Products</h1>
        <Link to="/products/new" className="btn">New Product</Link>
      </div>

      <form onSubmit={handleSearch} style={{ display: "flex", gap: 8, marginBottom: 16, maxWidth: "100%" }}>
        <input value={search} onChange={e => setSearch(e.target.value)} placeholder="Search by name, brand, barcode..." style={{ flex: 1 }} />
        <select value={categoryId} onChange={e => setSearchParams({ query, category_id: e.target.value, page: "1" })}>
          <option value="">All Categories</option>
          {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
        </select>
        <button type="submit">Search</button>
      </form>

      {loading ? <p>Loading...</p> : (
        <>
          <table>
            <thead>
              <tr><th>Name</th><th>Brand</th><th>Generic</th><th>Category</th><th>Strength</th><th>Form</th><th>Stock Alert</th><th>Actions</th></tr>
            </thead>
            <tbody>
              {products.map(p => (
                <tr key={p.id}>
                  <td><Link to={`/products/${p.id}`}>{p.name}</Link></td>
                  <td>{p.brand_name || "—"}</td>
                  <td>{p.generic_name || "—"}</td>
                  <td>{catMap[p.category_id] || "—"}</td>
                  <td>{p.strength || "—"}</td>
                  <td>{p.form || "—"}</td>
                  <td>{p.reorder_level}</td>
                  <td>
                    <Link to={`/products/${p.id}/edit`}>Edit</Link>
                    <button onClick={() => handleDelete(p.id)}>Delete</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {totalPages > 1 && (
            <div style={{ marginTop: 16, display: "flex", gap: 8 }}>
              {page > 1 && (
                <button onClick={() => setSearchParams({ query, category_id: categoryId, page: String(page - 1) })}>Previous</button>
              )}
              <span style={{ padding: "8px 0" }}>Page {page} of {totalPages}</span>
              {page < totalPages && (
                <button onClick={() => setSearchParams({ query, category_id: categoryId, page: String(page + 1) })}>Next</button>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
}
