# Stage 1: MotionGPT3
FROM pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime

# Set the working directory
WORKDIR /app

# Clone the MotionGPT3 repository
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y git && \
    git clone https://github.com/OpenMotionLab/MotionGPT3.git .

# Install dependencies
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    wget \
    ffmpeg \
    libgl1-mesa-glx \
    libgl1-mesa-dri \
    libglib2.0-0 \
    libsm6 \
    libxext6 \
    libxrender1 \
    libxi6 \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/davide-ferrara/animAI/raw/refs/heads/main/gpt3requirements.txt

# Install python dependencies, download models, and clean up
RUN pip install --no-cache-dir -r gpt3requirements.txt &&\
    ./prepare/download_smpl_model.sh && \
    ./prepare/prepare_gpt2.sh && \
    ./prepare/download_t2m_evaluators.sh && \
    ./prepare/download_mld_pretrained_models.sh && \
    ./prepare/download_pretrained_motiongpt3_model.sh && \
    python -m spacy download en_core_web_sm && \
    rm -rf /root/.cache && \
    rm -rf .git

# Copy the rest of the application
COPY . .
