package models

const (
	BLBase         = "https://api.beatleader.xyz"
	BLPlayerScores = "/player/%s/scores/compact?sortBy=pp&order=desc&count=100&page=%d"
	BLPlayers      = "/players?search=%s&leaderboardContext=general&page=1&count=10"

	SSBase         = "https://scoresaber.com/api"
	SSPlayers      = "/players?search=%s"
	SSPlayerScores = "/player/%s/scores?limit=60&withMetadata=true"
	SSPlayer       = "/api/player/%s/basic"

	BKGraphQlUrl       = "https://czqk28jt.apicdn.sanity.io/v1/graphql/prod_bk_de/default"
	BKCouponApiUrl     = "https://api.burgerking.de/api/o2uvrPdUY57J5WwYs6NtzZ2Knk7TnAUY/v4/de/de/coupons/"
	BKCouponImageUrl   = "https://cdn.sanity.io/images/czqk28jt/prod_bk_de/%s?w=512"
	BKCouponWebViewUrl = "https://www.burgerking.de/rewards/offers/%s"
)

var BKLoyaltyOffersQuery = `query MyQuery {
  allLoyaltyOffersUis {
    _id
    _updatedAt
	_type
    sortedSystemwideOffers {
      _id
      loyaltyEngineId
      shortCode
      offerPrice
      moreInfo {
        deRaw
      }
      vendorConfigs {
        partner {
          constantPlu
          discountPlu
          pluType
        }
      }
      rules {
        ... on LoyaltyBetweenDates {
          _key
          _type
          startDate
          endDate
        }
      }
      localizedImage {
        de {
          imageDescription
          app {
            asset {
              url
            }
          }
        }
      }
      name {
        deRaw
      }
    }
    _key
  }
}
`
var BKConfigOffersQuery = `query MyQuery {
  allConfigOffers {
    _createdAt
	_id
	_type
    internalName
    name {
      deRaw
    }
    localizedImage {
      de {
        imageDescription
        app {
          asset {
            url
          }
        }
      }
    }
    loyaltyEngineId
    moreInfo {
      deRaw
    }
    offerPrice
    redemptionType
    redemptionMethod
    shortCode
    vendorConfigs {
      partner {
        constantPlu
      }
    }
	description {
      deRaw
    }
  }
}
`
var BKSystemwideOffersQuery = `query MyQuery {
  allSystemwideOffers {
    _id
    _createdAt
    _updatedAt
    description {
      deRaw
    }
    internalName
    localizedImage {
      de {
        app {
          asset {
            _id
          }
        }
        imageDescription
      }
    }
    loyaltyEngineId
    moreInfo {
      deRaw
    }
    name {
      deRaw
    }
    offerPrice
    rules {
      ... on LoyaltyBetweenDates {
        _key
        _type
        endDate
        startDate
      }
    }
    shortCode
    vendorConfigs {
      partner {
        constantPlu
      }
    }
  }
}
`

type ResponseStructs interface {
	BLSongsResponse | BLPlayersResponse | SSPlayersResponse | ScoresResponseStructs
}
type ScoresResponseStructs interface {
	BLScoresResponse | SSScoresResponse
}
