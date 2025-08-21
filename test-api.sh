#!/bin/bash

# Test script for Stytch Member Management API
# This script demonstrates how to test the Lambda endpoints

# Configuration
API_ENDPOINT="https://srnext-stytch-member.sb.int.fullbayapi.com"
ORGANIZATION_ID="organization-test-00000000-0000-0000-0000-000000000000"  # Replace with actual org ID

echo "==================================================="
echo "Stytch Member Management API Test Script"
echo "==================================================="
echo ""

# Test health check endpoint
echo "1. Testing Health Check Endpoint..."
echo "   GET ${API_ENDPOINT}/members/health"
curl -s -X GET "${API_ENDPOINT}/members/health" | jq . || echo "Failed to connect"
echo ""

# Test create member
echo "2. Testing Create Member..."
echo "   POST ${API_ENDPOINT}/members"
MEMBER_DATA='{
  "organization_id": "'${ORGANIZATION_ID}'",
  "email": "test.user@example.com",
  "name": "Test User",
  "trusted_metadata": {
    "role": "user"
  }
}'

echo "   Request body:"
echo "${MEMBER_DATA}" | jq .
echo ""
echo "   Response:"
curl -s -X POST "${API_ENDPOINT}/members" \
  -H "Content-Type: application/json" \
  -d "${MEMBER_DATA}" | jq . || echo "Failed to create member"
echo ""

# Test get member (requires member_id from create response)
echo "3. Testing Get Member..."
echo "   GET ${API_ENDPOINT}/members/{member_id}?organization_id=${ORGANIZATION_ID}"
echo "   Note: Replace {member_id} with actual member ID from create response"
echo ""

# Test search members
echo "4. Testing Search Members..."
echo "   POST ${API_ENDPOINT}/members/search"
SEARCH_DATA='{
  "organization_id": "'${ORGANIZATION_ID}'",
  "limit": 10
}'

echo "   Request body:"
echo "${SEARCH_DATA}" | jq .
echo ""
echo "   Response:"
curl -s -X POST "${API_ENDPOINT}/members/search" \
  -H "Content-Type: application/json" \
  -d "${SEARCH_DATA}" | jq . || echo "Failed to search members"
echo ""

echo "==================================================="
echo "Test script completed"
echo "==================================================="