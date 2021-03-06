#Copyright (c) 2019 Web3 Technologies Foundation
#
#This file is part of Polkadot Host Test Suite
#
#Polkadot Host Test Suite is free software: you can redistribute it and/or modify
#it under the terms of the GNU General Public License as published by
#the Free Software Foundation, either version 3 of the License, or
#(at your option) any later version.
#
#Polkadot Host Tests is distributed in the hope that it will be useful,
#but WITHOUT ANY WARRANTY; without even the implied warranty of
#MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#GNU General Public License for more details.
#
#You should have received a copy of the GNU General Public License
#along with Foobar.  If not, see <https://www.gnu.org/licenses/>.

cmake_minimum_required(VERSION 3.12)


# Setup toolchain (required by kagome and libp2p)
set(CMAKE_TOOLCHAIN_FILE
  "${CMAKE_SOURCE_DIR}/cmake/Toolchain.cmake"
  CACHE FILEPATH "Default toolchain"
)


# Setup but disable hunter binary cache by default
set(
  HUNTER_USE_CACHE_SERVERS "NO"
  CACHE STRING "Binary cache server"
)

set(
  HUNTER_CACHE_SERVERS "https://github.com/soramitsu/hunter-binary-cache;"
  CACHE STRING "Binary cache server"
)


# Setup hunter
include(cmake/HunterGate.cmake)
HunterGate(
  URL https://github.com/soramitsu/soramitsu-hunter/archive/1a759980a8cfc0fcef847a293dbb4aceb081071d.zip 
  SHA1 4c8fd5dcd7b02b2cbb6f4d711da779032bdc8e82 
  FILEPATH "${CMAKE_SOURCE_DIR}/cmake/HunterConfig.cmake"
)

set(CMAKE_FIND_PACKAGE_PREFER_CONFIG ${HUNTER_ENABLED})


# General config
project(kagome-adapter-legacy LANGUAGES CXX C)


# Find dependencies
hunter_add_package(kagome)

find_package(kagome REQUIRED)

message(STATUS "Found kagome: ${kagome_INCLUDE_DIRS}")

# FIXME Kagome's package config should do all this!
find_package(schnorrkel_crust REQUIRED)
find_package(libsecp256k1 REQUIRED)
find_package(Boost REQUIRED COMPONENTS filesystem program_options random)
find_package(leveldb REQUIRED)
find_package(OpenSSL REQUIRED)
find_package(xxhash REQUIRED)
find_package(spdlog REQUIRED)
find_package(libp2p REQUIRED)
find_package(binaryen REQUIRED)


# Compile and link tester
add_executable(kagome-adapter-legacy
  src/main.cpp
  src/assert.cpp
  src/extension.cpp
  src/extension/storage.cpp
  src/extension/network.cpp
  src/extension/child_storage.cpp
  src/extension/crypto.cpp
  src/extension/helpers.cpp
)

# FIXME This should not be necessary, but not all public headers include from <kagome/..> (e.g. "outcome/outcome.hpp")
target_include_directories(kagome-adapter-legacy PRIVATE ${kagome_INCLUDE_DIRS})

target_link_libraries(kagome-adapter-legacy
  Boost::program_options
  ${kagome_LIBRARIES}
)


# Install result
include(GNUInstallDirs)

install(TARGETS kagome-adapter-legacy
  RUNTIME DESTINATION ${CMAKE_INSTALL_BINDIR}
)
