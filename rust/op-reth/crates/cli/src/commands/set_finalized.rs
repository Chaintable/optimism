//! Offline command for updating the persisted safe and finalized block numbers.

use alloy_primitives::{BlockHash, BlockNumber};
use clap::Parser;
use reth_db::{Database, open_db};
use reth_db_api::{
    tables::{self, ChainStateKey},
    transaction::{DbTx, DbTxMut},
};
use reth_node_core::args::DatabaseArgs;
use reth_stages::StageId;
use std::path::PathBuf;

/// Sets the persisted safe and finalized blocks in an offline op-reth database.
#[derive(Debug, Parser)]
pub struct Command {
    /// The op-reth data directory containing the `db` subdirectory.
    #[arg(long, value_name = "DATA_DIR")]
    datadir: PathBuf,

    /// All database related arguments.
    #[command(flatten)]
    db: DatabaseArgs,

    /// Block hash to mark safe and finalized. Defaults to the current local head.
    #[arg(value_name = "BLOCK_HASH")]
    block_hash: Option<BlockHash>,
}

impl Command {
    /// Executes the offline database update.
    pub fn execute(self) -> eyre::Result<()> {
        let db_path = self.datadir.join("db");
        eyre::ensure!(self.datadir.is_dir(), "Datadir does not exist: {:?}", self.datadir);
        eyre::ensure!(db_path.is_dir(), "Database does not exist: {db_path:?}");

        let db = open_db(&db_path, self.db.database_args())?;
        let result = set_finalized(&db, self.block_hash)?;

        println!(
            "Set safe and finalized block to #{}{} (previous safe: {}, previous finalized: {})",
            result.block_number,
            result.block_hash.map_or_else(String::new, |hash| format!(" ({hash})")),
            format_block_number(result.previous_safe),
            format_block_number(result.previous_finalized),
        );

        Ok(())
    }
}

#[derive(Debug, PartialEq, Eq)]
struct SetFinalizedResult {
    block_number: BlockNumber,
    block_hash: Option<BlockHash>,
    previous_safe: Option<BlockNumber>,
    previous_finalized: Option<BlockNumber>,
}

fn set_finalized<DB: Database>(
    db: &DB,
    block_hash: Option<BlockHash>,
) -> eyre::Result<SetFinalizedResult> {
    let tx = db.tx_mut()?;
    let block_number = match block_hash {
        Some(hash) => tx
            .get::<tables::HeaderNumbers>(hash)?
            .ok_or_else(|| eyre::eyre!("block hash not found: {hash}"))?,
        None => {
            tx.get::<tables::StageCheckpoints>(StageId::Finish.to_string())?
                .ok_or_else(|| {
                    eyre::eyre!("head block missing: Finish stage checkpoint not found")
                })?
                .block_number
        }
    };

    let previous_safe = tx.get::<tables::ChainState>(ChainStateKey::LastSafeBlock)?;
    let previous_finalized = tx.get::<tables::ChainState>(ChainStateKey::LastFinalizedBlock)?;

    tx.put::<tables::ChainState>(ChainStateKey::LastSafeBlock, block_number)?;
    tx.put::<tables::ChainState>(ChainStateKey::LastFinalizedBlock, block_number)?;
    tx.commit()?;

    Ok(SetFinalizedResult { block_number, block_hash, previous_safe, previous_finalized })
}

fn format_block_number(number: Option<BlockNumber>) -> String {
    number.map_or_else(|| "unset".to_string(), |number| format!("#{number}"))
}
