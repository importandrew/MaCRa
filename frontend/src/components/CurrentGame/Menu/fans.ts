// MCR fans in table order: number (some have a letter, e.g. "57 B"), name,
// description and points.
export type Fan = {
  no: string;
  name: string;
  description: string;
  points: number;
};

export const FANS: Fan[] = [
  {
    no: '1',
    name: 'Big Four Winds',
    description: 'Pungs or Kongs of of all four Wind Tiles.',
    points: 88,
  },
  {
    no: '2',
    name: 'Big Three Dragons',
    description: 'Pungs or Kongs of all three Dragon Tiles.',
    points: 88,
  },
  {
    no: '3',
    name: 'All Green',
    description:
      'A hand in which the chows, pungs and pair(s) are made up solely of “green” tiles: 2 Bam, 3 Bam, 4 Bam, 6 Bam, 8 Bam, and Green Dragon.',
    points: 88,
  },
  {
    no: '4',
    name: 'Nine Gates',
    description:
      'Holding the 1-1-1-2-3-4-5-6-7-8-9-9-9 tiles in any one of the suits, creating the nine-sided wait of 1-2-3-4-5-6-7-8-9. (Fully Concealed Hand may be combined if Self-Drawn).',
    points: 88,
  },
  {
    no: '5',
    name: 'Four Kongs',
    description:
      'Any hand that includes four kongs. They may be concealed or melded.',
    points: 88,
  },
  {
    no: '6',
    name: 'Seven Shifted Pairs',
    description:
      'A hand formed by seven pairs of the same suit, each shifted one up from the last. (Fully Concealed Hand may be combined if Self-Drawn).',
    points: 88,
  },
  {
    no: '7',
    name: 'Thirteen Orphans',
    description:
      'A hand created by singles of any 12 of the 1, 9, and Honor tiles, along with a pair of the 13th. (Fully Concealed Hand may be combined if Self-Drawn).',
    points: 88,
  },
  {
    no: '8',
    name: 'All Terminals',
    description:
      'The pair(s), Pungs or Kongs are all made up of 1 or 9 Number Tiles, without Honor Tiles.',
    points: 64,
  },
  {
    no: '9',
    name: 'Little Four Winds',
    description:
      'A hand that includes three Pungs or Kongs of Wind Tiles, and a pair of the fourth Wind.',
    points: 64,
  },
  {
    no: '10',
    name: 'Little Three Dragons',
    description:
      'A hand that includes two Pungs or Kongs of the Dragon Tiles, and a pair of the third Dragon.',
    points: 64,
  },
  {
    no: '11',
    name: 'All Honors',
    description: 'The pair(s), Pungs or Kongs are all made up of Honor Tiles.',
    points: 64,
  },
  {
    no: '12',
    name: 'Four Concealed Pungs',
    description:
      'A hand that includes four Concealed Pungs or Kongs (achieved without melding). (Fully Concealed Hand may be combined if Self-Drawn).',
    points: 64,
  },
  {
    no: '13',
    name: 'Pure Terminal Chows',
    description:
      'A hand consisting of two each of the lower and upper terminal Chows in one suit only, and a pair of fives in the same suit.',
    points: 64,
  },
  {
    no: '14',
    name: 'Quadruple Chow',
    description:
      'Four chows of the same continuous number sequence in the same suit.',
    points: 48,
  },
  {
    no: '15',
    name: 'Four Pure Shifted Pungs',
    description:
      'Four Pungs (or Kongs) in the same suit, each shifted one up from the last.',
    points: 48,
  },
  {
    no: '16',
    name: 'Four Pure Shifted Chows',
    description:
      'Four chows in one suit, each shifted up 1 or 2 numbers from the last, but not a combination of both.',
    points: 32,
  },
  {
    no: '17',
    name: 'Three Kongs',
    description:
      'A hand containing three Kongs. (Points for concealment may be added).',
    points: 32,
  },
  {
    no: '18',
    name: 'All Terminals and Honors',
    description:
      'The Pair(s), Pungs or Kongs are all made up of 1 or 9 Number Tiles and Honor Tiles.',
    points: 32,
  },
  {
    no: '19',
    name: 'Seven Pairs',
    description:
      'A Hand formed by seven Pairs. (Fully Concealed Hand may be combined if Self-Drawn).',
    points: 24,
  },
  {
    no: '20',
    name: 'Greater Honors and Knitted Tiles',
    description:
      'Formed by seven single honors and singles of suit tiles belonging to separate Knitted sequences (for example, 1-4-7 of Bamboos, 2-5-8 of Characters, and 3-6-9 of Dots). Fully Concealed Hand may be combined if Self-Drawn.',
    points: 24,
  },
  {
    no: '21',
    name: 'All Even Pungs',
    description:
      'A Hand formed with Pungs or Kongs of 2, 4, 6, and 8 tiles, with a pair of the same.',
    points: 24,
  },
  {
    no: '22',
    name: 'Full Flush',
    description: 'A Hand formed entirely of a single suit.',
    points: 24,
  },
  {
    no: '23',
    name: 'Pure Triple Chow',
    description:
      'Three chows of the same numerical sequence and in the same suit.',
    points: 24,
  },
  {
    no: '24',
    name: 'Pure Shifted Pungs',
    description:
      'Three Pungs or Kongs of the same suit, each shifted one up from the last.',
    points: 24,
  },
  {
    no: '25',
    name: 'Upper Tiles',
    description: 'A Hand consisting entirely of 7, 8, and 9 tiles.',
    points: 24,
  },
  {
    no: '26',
    name: 'Middle Tiles',
    description: 'A Hand consisting entirely of 4, 5, and 6 tiles.',
    points: 24,
  },
  {
    no: '27',
    name: 'Lower Tiles',
    description: 'A Hand consisting entirely of 1, 2, and 3 tiles.',
    points: 24,
  },
  {
    no: '28',
    name: 'Pure Straight',
    description:
      'A Hand using one each of all the numbers 1 through 9 from any one suit, forming three consecutive chows.',
    points: 16,
  },
  {
    no: '29',
    name: 'Three-Suited Terminal Chows',
    description:
      'A Hand consisting of 1-2-3 and 7-8-9 in one suit (Two Terminal Chows), 1-2-3 and 7-8-9 in another suit, and a pair of fives in the remaining suit.',
    points: 16,
  },
  {
    no: '30',
    name: 'Pure Shifted Chows',
    description:
      'Three Chows in one suit, each shifted either one or two numbers up from the last, but not a combination of both.',
    points: 16,
  },
  {
    no: '31',
    name: 'All Fives',
    description:
      'A Hand in which every set (chow, pung, kong, pair) includes the number “5”.',
    points: 16,
  },
  {
    no: '32',
    name: 'Triple Pung',
    description: 'Three Pungs (or Kongs) of the same number in each suit.',
    points: 16,
  },
  {
    no: '33',
    name: 'Three Concealed Pungs',
    description: 'Three Concealed Pungs or Kongs (achieved without melding).',
    points: 16,
  },
  {
    no: '34',
    name: 'Lesser Honors and Knitted Tiles',
    description:
      'Formed by single honors, and singles of suit tiles belonging to separate Knitted sequences (for example, 1-4-7 of Bamboo, 2-5-8 of Characters, and 3-6-9 of Dots - each of the 3 suits must belong to a different Knitted sequence, but not necessarily in this order). Fully Concealed Hand may be combined if Self-Drawn.',
    points: 12,
  },
  {
    no: '35',
    name: 'Knitted Straight',
    description:
      'A special Straight which is formed not with standard chows but with 3 different Knitted sequences. For example, 1-4-7 of Dots, 2-5-8 of Characters, and 3-6-9 of Bamboos - but not necessarily in the order in this example.',
    points: 12,
  },
  {
    no: '36',
    name: 'Upper Four',
    description: 'A Hand created with suit tiles 6 through 9.',
    points: 12,
  },
  {
    no: '37',
    name: 'Lower Four',
    description: 'A Hand created with suit tiles 1 through 4.',
    points: 12,
  },
  {
    no: '38',
    name: 'Big Three Winds',
    description:
      'A Hand that includes one pung (or kong) of each of the three winds.',
    points: 12,
  },
  {
    no: '39',
    name: 'Mixed Straight',
    description:
      'A Straight (tiles 1 through 9) formed by chows from all three suits.',
    points: 8,
  },
  {
    no: '40',
    name: 'Reversible Tiles',
    description:
      'A Hand created entirely with those tiles which are vertically symmetrical, which means the carved designs look the same if you turn them upside-down. These tiles are the 1,2,3,4,5,8, and 9 Dots, the 2,4,5,6,8, and 9 Bams, and the White Dragon.',
    points: 8,
  },
  {
    no: '41',
    name: 'Mixed Triple Chow',
    description:
      'Three Chows of the same numerical sequence, one in each suit.',
    points: 8,
  },
  {
    no: '42',
    name: 'Mixed Shifted Pungs',
    description:
      'Three Pungs (or kongs), one in each suit, each shifted up one number from the last.',
    points: 8,
  },
  {
    no: '43',
    name: 'Chicken Hand',
    description:
      'A Hand that would otherwise earn 0 points (excluding the Flower Tiles).',
    points: 8,
  },
  {
    no: '44',
    name: 'Last Tile Draw',
    description:
      'Going out (making Mahjong) on a pick of the very last tile of the wall. (Points for Self-Drawn may not be combined.)',
    points: 8,
  },
  {
    no: '45',
    name: 'Last Tile Claim',
    description:
      'Going out (making Mahjong) off the discard which is the last tile in the game.',
    points: 8,
  },
  {
    no: '46',
    name: 'Out with Replacement Tile',
    description:
      'Going out (making mahjong) on the replacement tile drawn after achieving a kong (not on a Flower replacement). When a Flower Tile is taken after Konging, and upon winning on the Flower replacement, points for Self-Drawn may be added (but Out With Replacement Tile does not apply in this case).',
    points: 8,
  },
  {
    no: '47',
    name: 'Robbing The Kong',
    description:
      'Winning off the tile that somebody adds to a melded pung (to create a Kong). (The points for Last Tile may not be combined.)',
    points: 8,
  },
  {
    no: '48',
    name: 'Two Concealed Kongs',
    description: 'A Hand that includes two Concealed Kongs.',
    points: 8,
  },
  {
    no: '49',
    name: 'All Pungs',
    description: 'A Hand formed by four Pungs (or Kongs) and one pair.',
    points: 6,
  },
  {
    no: '50',
    name: 'Half Flush',
    description:
      'A Hand formed by tiles from any one of the three suits, in combination with Honor tiles.',
    points: 6,
  },
  {
    no: '51',
    name: 'Mixed Shifted Chows',
    description:
      'Three Chows, one in each suit, each shifted up one number from the last.',
    points: 6,
  },
  {
    no: '52',
    name: 'All Types',
    description:
      'A Hand in which each of the five sets (pungs, kongs, chows, pairs) is composed of a different type of tile (Characters, Bamboos, Dots, Winds, and Dragons).',
    points: 6,
  },
  {
    no: '53',
    name: 'Melded Hand',
    description:
      'Every set in the hand (chow, pung, kong, and pair) must be completed with tiles discarded by other players. All sets must be exposed, and the player goes out on a single wait off another player.',
    points: 6,
  },
  {
    no: '54',
    name: 'Two Dragon Pungs',
    description: 'Two Pungs (or kongs) of Dragon tiles.',
    points: 6,
  },
  {
    no: '57 B',
    name: 'One Melded and one Concealed Kongs',
    description: 'A Hand that includes one Melded and one Concealed Kongs.',
    points: 6,
  },
  {
    no: '55',
    name: 'Outside Hand',
    description:
      'A Hand that includes terminals and honors in each set, including the pair.',
    points: 4,
  },
  {
    no: '56',
    name: 'Fully Concealed Hand',
    description:
      'A Hand that a player completes without any melds, and wins by Self-Draw.',
    points: 4,
  },
  {
    no: '57 A',
    name: 'Two Melded Kongs',
    description: 'A Hand that includes two Melded Kongs.',
    points: 4,
  },
  {
    no: '58',
    name: 'Last Tile',
    description:
      'Winning on a tile that is the last of its kind. (It must be clear to all players based on the discards and exposures.)',
    points: 4,
  },
  {
    no: '59',
    name: 'Dragon Pung',
    description: 'A Pung or Kong of Dragon Tiles.',
    points: 2,
  },
  {
    no: '60',
    name: 'Prevalent Wind',
    description:
      'A Pung or Kong of the Wind Tile corresponding to the current Prevalent Wind.',
    points: 2,
  },
  {
    no: '61',
    name: 'Seat Wind',
    description:
      'A Pung or Kong of the Wind Tile corresponding to the player’s Seat position at the table. (Dealer is East; proceeding counter-clockwise from the Dealer, other players’ seats are South, West, North.)',
    points: 2,
  },
  {
    no: '62',
    name: 'Concealed Hand',
    description:
      'Having a concealed Hand (no melded sets) and winning by discard.',
    points: 2,
  },
  {
    no: '63',
    name: 'All Chows',
    description: 'A hand consisting of all chows, with no Honors.',
    points: 2,
  },
  {
    no: '64',
    name: 'Tile Hog',
    description:
      'Using all four of a single suit tile, without using them as a Kong.',
    points: 2,
  },
  {
    no: '65',
    name: 'Double Pung',
    description:
      'Two Pungs (or Kongs) of the same number in two different suits.',
    points: 2,
  },
  {
    no: '66',
    name: 'Two Concealed Pungs',
    description: 'Two Pungs achieved without melding.',
    points: 2,
  },
  {
    no: '67',
    name: 'Concealed Kong',
    description:
      'Created when four identical tiles, all self-drawn, are declared as a Kong.',
    points: 2,
  },
  {
    no: '68',
    name: 'All Simples',
    description: 'A Hand formed without Terminal or Honor Tiles.',
    points: 2,
  },
  {
    no: '69',
    name: 'Pure Double Chow',
    description: 'Two identical chows in the same suit.',
    points: 1,
  },
  {
    no: '70',
    name: 'Mixed Double Chow',
    description: 'Two chows of the same numbers but in different suits.',
    points: 1,
  },
  {
    no: '71',
    name: 'Short Straight',
    description:
      'Two chows in the same suit that run consecutively after one another to make a six-tile straight.',
    points: 1,
  },
  {
    no: '72',
    name: 'Two Terminal Chows',
    description: 'Chows of 1-2-3 and 7-8-9 in the same suit.',
    points: 1,
  },
  {
    no: '73',
    name: 'Pung of Terminals or Honors',
    description:
      'A Pung or Kong of Ones, Nines, or Winds. (A dragon pung scores 2 points.)',
    points: 1,
  },
  {
    no: '74',
    name: 'Melded Kong',
    description:
      'A Kong that was claimed from another player or promoted from a melded pung.',
    points: 1,
  },
  {
    no: '75',
    name: 'One Voided Suit',
    description:
      'A Hand that uses tiles from only two of the three suits (it lacks any tiles from one of the three suits).',
    points: 1,
  },
  {
    no: '76',
    name: 'No Honors',
    description:
      'A Hand formed entirely of suit tiles, without Winds or Dragons.',
    points: 1,
  },
  {
    no: '77',
    name: 'Edge Wait',
    description:
      'Waiting solely for a 3 to form a 1-2-3 chow, or solely for a 7 to form a 7-8-9 chow. Not valid if waiting for more than one tile. Not valid if the edge wait is combined with any other waits.',
    points: 1,
  },
  {
    no: '78',
    name: 'Closed Wait',
    description:
      'Waiting solely for a tile whose number is “inside” (in the middle) to form a chow. Not valid if waiting for more than one tile. Not valid if the closed wait is combined with other waits.',
    points: 1,
  },
  {
    no: '79',
    name: 'Single Wait',
    description:
      'Waiting solely for a tile to form a pair. Not valid if waiting for more than one tile (for example, holding 1-2-3-4 and waiting on the 1 and 4).',
    points: 1,
  },
  {
    no: '80',
    name: 'Self-Drawn',
    description:
      'Going out (making mahjong) with a fresh tile picked from the wall.',
    points: 1,
  },
  {
    no: '81',
    name: 'Flower Tiles',
    description:
      'Each tile carved with Chinese word of Spring (or Summer, Autumn, Winder, Plum, Orchid, Bamboo, Chrysanthemum) will award you one point when you succeed in Hu. Flower replacement to be Hu, the point of Self-Drawn can add. It can’t add the points of Hu While Kong, it not Flower replacement can discard.',
    points: 1,
  },
];
