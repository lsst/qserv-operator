"""Sphinx configurations for the qserv.lsst.io documentation build."""

import os
import sys

import lsst_sphinx_bootstrap_theme

# -- General configuration ----------------------------------------------------

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
extensions = [
    "sphinx.ext.intersphinx",
    "sphinx.ext.ifconfig",
    "documenteer.sphinxext",
]

source_suffix = ".rst"

root_doc = "index"
master_doc = "index"  # deprecated in 4.0

project = "Qserv-operator"

# TODO
copyright = ""

author = "Rubin Observatory"

# The version info for the project you're documenting, acts as replacement for
# |version| and |release|, also used in various other places throughout the
# built documents.
# TODO either set these manually here, or extract from e.g. env variables.
version = "X.Y.Z"
release = "X.Y.Z"

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
exclude_patterns = ["_build", "Makefile"]

# The reST default role cross-links Python (used for this markup: `text`)
default_role = "py:obj"

# -- ReStructuredText epilog for common links/substitutions ---------------

rst_epilog = """

.. _conda-forge: https://conda-forge.org
.. _conda: https://conda.io/en/latest/index.html
"""

# -- Options for linkcheck builder --------------------------------------------

linkcheck_retries = 2

linkcheck_timeout = 15

# Add any URL patterns to ignore (e.g. for private sites, or sites that
# are frequently down).
linkcheck_ignore = [
    r"^https://jira.lsstcorp.org/browse/",
    r"^https://dev.lsstcorp.org/trac"
]

# -- Options for html builder -------------------------------------------------

templates_path = [
    lsst_sphinx_bootstrap_theme.get_html_templates_path(),
]
html_theme = "lsst_sphinx_bootstrap_theme"
html_theme_path = [lsst_sphinx_bootstrap_theme.get_html_theme_path()]

# Variables available for Jinja templates
html_context = {}

# Theme options are theme-specific and customize the look and feel of a theme
# further.  For a list of options available for each theme, see the
# documentation.
html_theme_options = {"logotext": project}

# The name for this set of Sphinx documents.  If unset, it defaults to
# "<project> v<release> documentation".
# html_title = ""

# A shorter title for the navigation bar.  Default is the same as html_title.
html_short_title = "Qserv-operator"

# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ["_static"]

# If true, links to the reST sources are added to the pages.
html_show_sourcelink = False

# -- Intersphinx --------------------------------------------------------------
# For linking to other Sphinx documentation.
# https://www.sphinx-doc.org/en/master/usage/extensions/intersphinx.html

intersphinx_mapping = {
    "python": ("https://docs.python.org/3/", None),
    "pipelines": ("https://pipelines.lsst.io/", None),
}
