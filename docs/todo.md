- `[2026-07-26][decided] Keep proof-history trie on the same Reth source` — Alethia consumes `reth-optimism-trie`, whose workspace pins the same Reth revision as Alethia.
  **Decision:** Pin the Optimism Rust workspace to `DeBankDeFi/reth@f077b171e5cbcfe315d51c9fb71d042d87c9af07` so all exposed provider/trie types come from one Cargo source. Do not allow duplicate paradigmxyz/DeBankDeFi Reth packages in the Alethia graph.

- `[2026-07-26][done] Dependency fork verification` — Refresh the Rust lockfile and run the focused `reth-optimism-trie` checks before publishing the exact Optimism commit for Alethia.
  **Done:** `cargo check -p reth-optimism-trie` and the repository-pinned nightly rustfmt check passed with all 105 Reth packages resolved from `DeBankDeFi/reth@f077b171`. Focused clippy passed after allowing three pre-existing `useless_conversion` / `missing_const_for_fn` findings in unchanged trie source.
