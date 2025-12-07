import argparse
import logging
from pipeline.runner import PipelineRunner

logging.basicConfig(
    level=logging.INFO, # Only display messages of INFO level and above
    format='%(asctime)s - %(levelname)s - %(message)s'
)

logger = logging.getLogger(__name__)

def _parse_cli_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Manga Data Pipeline CLI.")   
    parser.add_argument(
        '--init', 
        action='store_true', # This makes the argument a simple True/False flag
        help='Run the top 100 initialization pipeline.'
    )
    
    return parser.parse_args()

def main():
    arg = _parse_cli_args()
    runner = PipelineRunner()

    if arg.init:
        logger.info("Starting Top 100 Manga Initialization Pipeline...")
        runner.run_top_100_init_pipline()
    else:
        logger.info("Staring Update Lastest Chapter Pipeline...")
        runner.run_update_non_complete_pipeline()

if __name__ == "__main__":
    main()