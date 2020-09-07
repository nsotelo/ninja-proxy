from setuptools import setup

setup(
    author="Nigel Sotelo",
    author_email="nigelrsotelo@gmail.com",
    description="Generate links for ninja-proxy",
    entry_points={"console_scripts": ["ninja-link=utils.generate_link:main", "ninja-key=utils.generate_key:main"]},
    install_requires=["pycryptodome"],
    name="ninja-proxy-utils",
    packages=["utils"],
    version="1.0",
)
