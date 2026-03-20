import {
  Link,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
} from '@heroui/react';

export function FansMenu() {
  return (
    <div
      className={
        'overflow-hidden flex flex-col justify-center mb-3 w-full ' +
        '[&_table]:w-full [&_table]:table-fixed [&_:is(th,td)]:box-border ' +
        '[&_:is(th,td):first-child]:w-10 [&_:is(th,td):first-child]:min-w-10 ' +
        '[&_:is(th,td):last-child]:w-10 [&_:is(th,td):last-child]:min-w-10 ' +
        '[&_:is(th,td):nth-child(n+2):nth-child(-n+2)]:w-[calc((100%-5rem)/3)] [&_:is(th,td):nth-child(n+2):nth-child(-n+2)]:min-w-0'
      }
    >
      <Table
        isStriped
        removeWrapper
        color="default"
        selectionMode="single"
        aria-label="Fans list"
        className="sm:w-fit self-center"
      >
        <TableHeader>
          <TableColumn>
            <strong>No.</strong>
          </TableColumn>
          <TableColumn>
            <strong>Fan name</strong>
          </TableColumn>
          <TableColumn>
            <strong>Description</strong>
          </TableColumn>
          <TableColumn>
            <strong>Pts.</strong>
          </TableColumn>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableCell>1</TableCell>
            <TableCell>Big Four Winds</TableCell>
            <TableCell>Pungs or Kongs of of all four Wind Tiles.</TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>2</TableCell>
            <TableCell>Big Three Dragons</TableCell>
            <TableCell>Pungs or Kongs of all three Dragon Tiles.</TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>3</TableCell>
            <TableCell>All Green</TableCell>
            <TableCell>
              A hand in which the chows, pungs and pair(s) are made up solely of
              “green” tiles: 2 Bam, 3 Bam, 4 Bam, 6 Bam, 8 Bam, and Green
              Dragon.
            </TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>4</TableCell>
            <TableCell>Nine Gates</TableCell>
            <TableCell>
              Holding the 1-1-1-2-3-4-5-6-7-8-9-9-9 tiles in any one of the
              suits, creating the nine-sided wait of 1-2-3-4-5-6-7-8-9. (Fully
              Concealed Hand may be combined if Self-Drawn).
            </TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>5</TableCell>
            <TableCell>Four Kongs</TableCell>
            <TableCell>
              Any hand that includes four kongs. They may be concealed or
              melded.
            </TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>6</TableCell>
            <TableCell>Seven Shifted Pairs</TableCell>
            <TableCell>
              A hand formed by seven pairs of the same suit, each shifted one up
              from the last. (Fully Concealed Hand may be combined if
              Self-Drawn).
            </TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>7</TableCell>
            <TableCell>Thirteen Orphans</TableCell>
            <TableCell>
              A hand created by singles of any 12 of the 1, 9, and Honor tiles,
              along with a pair of the 13th. (Fully Concealed Hand may be
              combined if Self-Drawn).
            </TableCell>
            <TableCell>88</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>8</TableCell>
            <TableCell>All Terminals</TableCell>
            <TableCell>
              The pair(s), Pungs or Kongs are all made up of 1 or 9 Number
              Tiles, without Honor Tiles.
            </TableCell>
            <TableCell>64</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>9</TableCell>
            <TableCell>Little Four Winds</TableCell>
            <TableCell>
              A hand that includes three Pungs or Kongs of Wind Tiles, and a
              pair of the fourth Wind.
            </TableCell>
            <TableCell>64</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>10</TableCell>
            <TableCell>Little Three Dragons</TableCell>
            <TableCell>
              A hand that includes two Pungs or Kongs of the Dragon Tiles, and a
              pair of the third Dragon.
            </TableCell>
            <TableCell>64</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>11</TableCell>
            <TableCell>All Honors</TableCell>
            <TableCell>
              The pair(s), Pungs or Kongs are all made up of Honor Tiles.
            </TableCell>
            <TableCell>64</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>12</TableCell>
            <TableCell>Four Concealed Pungs</TableCell>
            <TableCell>
              A hand that includes four Concealed Pungs or Kongs (achieved
              without melding). (Fully Concealed Hand may be combined if
              Self-Drawn).
            </TableCell>
            <TableCell>64</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>13</TableCell>
            <TableCell>Pure Terminal Chows</TableCell>
            <TableCell>
              A hand consisting of two each of the lower and upper terminal
              Chows in one suit only, and a pair of fives in the same suit.
            </TableCell>
            <TableCell>64</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>14</TableCell>
            <TableCell>Quadruple Chow</TableCell>
            <TableCell>
              Four chows of the same continuous number sequence in the same
              suit.
            </TableCell>
            <TableCell>48</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>15</TableCell>
            <TableCell>Four Pure Shifted Pungs</TableCell>
            <TableCell>
              Four Pungs (or Kongs) in the same suit, each shifted one up from
              the last.
            </TableCell>
            <TableCell>48</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>16</TableCell>
            <TableCell>Four Pure Shifted Chows</TableCell>
            <TableCell>
              Four chows in one suit, each shifted up 1 or 2 numbers from the
              last, but not a combination of both.
            </TableCell>
            <TableCell>32</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>17</TableCell>
            <TableCell>Three Kongs</TableCell>
            <TableCell>
              A hand containing three Kongs. (Points for concealment may be
              added).
            </TableCell>
            <TableCell>32</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>18</TableCell>
            <TableCell>All Terminals and Honors</TableCell>
            <TableCell>
              The Pair(s), Pungs or Kongs are all made up of 1 or 9 Number Tiles
              and Honor Tiles.
            </TableCell>
            <TableCell>32</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>19</TableCell>
            <TableCell>Seven Pairs</TableCell>
            <TableCell>
              A Hand formed by seven Pairs. (Fully Concealed Hand may be
              combined if Self-Drawn).
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>20</TableCell>
            <TableCell>Greater Honors and Knitted Tiles</TableCell>
            <TableCell>
              Formed by seven single honors and singles of suit tiles belonging
              to separate Knitted sequences (for example, 1-4-7 of Bamboos,
              2-5-8 of Characters, and 3-6-9 of Dots). Fully Concealed Hand may
              be combined if Self-Drawn.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>21</TableCell>
            <TableCell>All Even Pungs</TableCell>
            <TableCell>
              A Hand formed with Pungs or Kongs of 2, 4, 6, and 8 tiles, with a
              pair of the same.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>22</TableCell>
            <TableCell>Full Flush</TableCell>
            <TableCell>A Hand formed entirely of a single suit.</TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>23</TableCell>
            <TableCell>Pure Triple Chow</TableCell>
            <TableCell>
              Three chows of the same numerical sequence and in the same suit.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>24</TableCell>
            <TableCell>Pure Shifted Pungs</TableCell>
            <TableCell>
              Three Pungs or Kongs of the same suit, each shifted one up from
              the last.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>25</TableCell>
            <TableCell>Upper Tiles</TableCell>
            <TableCell>
              A Hand consisting entirely of 7, 8, and 9 tiles.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>26</TableCell>
            <TableCell>Middle Tiles</TableCell>
            <TableCell>
              A Hand consisting entirely of 4, 5, and 6 tiles.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>27</TableCell>
            <TableCell>Lower Tiles</TableCell>
            <TableCell>
              A Hand consisting entirely of 1, 2, and 3 tiles.
            </TableCell>
            <TableCell>24</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>28</TableCell>
            <TableCell>Pure Straight</TableCell>
            <TableCell>
              A Hand using one each of all the numbers 1 through 9 from any one
              suit, forming three consecutive chows.
            </TableCell>
            <TableCell>16</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>29</TableCell>
            <TableCell>Three-Suited Terminal Chows</TableCell>
            <TableCell>
              A Hand consisting of 1-2-3 and 7-8-9 in one suit (Two Terminal
              Chows), 1-2-3 and 7-8-9 in another suit, and a pair of fives in
              the remaining suit.
            </TableCell>
            <TableCell>16</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>30</TableCell>
            <TableCell>Pure Shifted Chows</TableCell>
            <TableCell>
              Three Chows in one suit, each shifted either one or two numbers up
              from the last, but not a combination of both.
            </TableCell>
            <TableCell>16</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>31</TableCell>
            <TableCell>All Fives</TableCell>
            <TableCell>
              A Hand in which every set (chow, pung, kong, pair) includes the
              number “5”.
            </TableCell>
            <TableCell>16</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>32</TableCell>
            <TableCell>Triple Pung</TableCell>
            <TableCell>
              Three Pungs (or Kongs) of the same number in each suit.
            </TableCell>
            <TableCell>16</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>33</TableCell>
            <TableCell>Three Concealed Pungs</TableCell>
            <TableCell>
              Three Concealed Pungs or Kongs (achieved without melding).
            </TableCell>
            <TableCell>16</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>34</TableCell>
            <TableCell>Lesser Honors and Knitted Tiles</TableCell>
            <TableCell>
              Formed by single honors, and singles of suit tiles belonging to
              separate Knitted sequences (for example, 1-4-7 of Bamboo, 2-5-8 of
              Characters, and 3-6-9 of Dots - each of the 3 suits must belong to
              a different Knitted sequence, but not necessarily in this order).
              Fully Concealed Hand may be combined if Self-Drawn.
            </TableCell>
            <TableCell>12</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>35</TableCell>
            <TableCell>Knitted Straight</TableCell>
            <TableCell>
              A special Straight which is formed not with standard chows but
              with 3 different Knitted sequences. For example, 1-4-7 of Dots,
              2-5-8 of Characters, and 3-6-9 of Bamboos - but not necessarily in
              the order in this example.
            </TableCell>
            <TableCell>12</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>36</TableCell>
            <TableCell>Upper Four</TableCell>
            <TableCell>A Hand created with suit tiles 6 through 9.</TableCell>
            <TableCell>12</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>37</TableCell>
            <TableCell>Lower Four</TableCell>
            <TableCell>A Hand created with suit tiles 1 through 4.</TableCell>
            <TableCell>12</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>38</TableCell>
            <TableCell>Big Three Winds</TableCell>
            <TableCell>
              A Hand that includes one pung (or kong) of each of the three
              winds.
            </TableCell>
            <TableCell>12</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>39</TableCell>
            <TableCell>Mixed Straight</TableCell>
            <TableCell>
              A Straight (tiles 1 through 9) formed by chows from all three
              suits.
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>40</TableCell>
            <TableCell>Reversible Tiles</TableCell>
            <TableCell>
              A Hand created entirely with those tiles which are vertically
              symmetrical, which means the carved designs look the same if you
              turn them upside-down. These tiles are the 1,2,3,4,5,8, and 9
              Dots, the 2,4,5,6,8, and 9 Bams, and the White Dragon.
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>41</TableCell>
            <TableCell>Mixed Triple Chow</TableCell>
            <TableCell>
              Three Chows of the same numerical sequence, one in each suit.
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>42</TableCell>
            <TableCell>Mixed Shifted Pungs</TableCell>
            <TableCell>
              Three Pungs (or kongs), one in each suit, each shifted up one
              number from the last.
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>43</TableCell>
            <TableCell>Chicken Hand</TableCell>
            <TableCell>
              A Hand that would otherwise earn 0 points (excluding the Flower
              Tiles).
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>44</TableCell>
            <TableCell>Last Tile Draw</TableCell>
            <TableCell>
              Going out (making Mahjong) on a pick of the very last tile of the
              wall. (Points for Self-Drawn may not be combined.)
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>45</TableCell>
            <TableCell>Last Tile Claim</TableCell>
            <TableCell>
              Going out (making Mahjong) off the discard which is the last tile
              in the game.
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>46</TableCell>
            <TableCell>Out with Replacement Tile</TableCell>
            <TableCell>
              Going out (making mahjong) on the replacement tile drawn after
              achieving a kong (not on a Flower replacement). When a Flower Tile
              is taken after Konging, and upon winning on the Flower
              replacement, points for Self-Drawn may be added (but Out With
              Replacement Tile does not apply in this case).
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>47</TableCell>
            <TableCell>Robbing The Kong</TableCell>
            <TableCell>
              Winning off the tile that somebody adds to a melded pung (to
              create a Kong). (The points for Last Tile may not be combined.)
            </TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>48</TableCell>
            <TableCell>Two Concealed Kongs</TableCell>
            <TableCell>A Hand that includes two Concealed Kongs.</TableCell>
            <TableCell>8</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>49</TableCell>
            <TableCell>All Pungs</TableCell>
            <TableCell>
              A Hand formed by four Pungs (or Kongs) and one pair.
            </TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>50</TableCell>
            <TableCell>Half Flush</TableCell>
            <TableCell>
              A Hand formed by tiles from any one of the three suits, in
              combination with Honor tiles.
            </TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>51</TableCell>
            <TableCell>Mixed Shifted Chows</TableCell>
            <TableCell>
              Three Chows, one in each suit, each shifted up one number from the
              last.
            </TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>52</TableCell>
            <TableCell>All Types</TableCell>
            <TableCell>
              A Hand in which each of the five sets (pungs, kongs, chows, pairs)
              is composed of a different type of tile (Characters, Bamboos,
              Dots, Winds, and Dragons).
            </TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>53</TableCell>
            <TableCell>Melded Hand</TableCell>
            <TableCell>
              Every set in the hand (chow, pung, kong, and pair) must be
              completed with tiles discarded by other players. All sets must be
              exposed, and the player goes out on a single wait off another
              player.
            </TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>54</TableCell>
            <TableCell>Two Dragon Pungs</TableCell>
            <TableCell>Two Pungs (or kongs) of Dragon tiles.</TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>57 B</TableCell>
            <TableCell>One Melded and one Concealed Kongs</TableCell>
            <TableCell>
              A Hand that includes one Melded and one Concealed Kongs.
            </TableCell>
            <TableCell>6</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>55</TableCell>
            <TableCell>Outside Hand</TableCell>
            <TableCell>
              A Hand that includes terminals and honors in each set, including
              the pair.
            </TableCell>
            <TableCell>4</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>56</TableCell>
            <TableCell>Fully Concealed Hand</TableCell>
            <TableCell>
              A Hand that a player completes without any melds, and wins by
              Self-Draw.
            </TableCell>
            <TableCell>4</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>57 A</TableCell>
            <TableCell>Two Melded Kongs</TableCell>
            <TableCell>A Hand that includes two Melded Kongs.</TableCell>
            <TableCell>4</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>58</TableCell>
            <TableCell>Last Tile</TableCell>
            <TableCell>
              Winning on a tile that is the last of its kind. (It must be clear
              to all players based on the discards and exposures.)
            </TableCell>
            <TableCell>4</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>59</TableCell>
            <TableCell>Dragon Pung</TableCell>
            <TableCell>A Pung or Kong of Dragon Tiles.</TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>60</TableCell>
            <TableCell>Prevalent Wind</TableCell>
            <TableCell>
              A Pung or Kong of the Wind Tile corresponding to the current
              Prevalent Wind.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>61</TableCell>
            <TableCell>Seat Wind</TableCell>
            <TableCell>
              A Pung or Kong of the Wind Tile corresponding to the player’s Seat
              position at the table. (Dealer is East; proceeding
              counter-clockwise from the Dealer, other players’ seats are South,
              West, North.)
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>62</TableCell>
            <TableCell>Concealed Hand</TableCell>
            <TableCell>
              Having a concealed Hand (no melded sets) and winning by discard.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>63</TableCell>
            <TableCell>All Chows</TableCell>
            <TableCell>
              A hand consisting of all chows, with no Honors.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>64</TableCell>
            <TableCell>Tile Hog</TableCell>
            <TableCell>
              Using all four of a single suit tile, without using them as a
              Kong.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>65</TableCell>
            <TableCell>Double Pung</TableCell>
            <TableCell>
              Two Pungs (or Kongs) of the same number in two different suits.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>66</TableCell>
            <TableCell>Two Concealed Pungs</TableCell>
            <TableCell>Two Pungs achieved without melding.</TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>67</TableCell>
            <TableCell>Concealed Kong</TableCell>
            <TableCell>
              Created when four identical tiles, all self-drawn, are declared as
              a Kong.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>68</TableCell>
            <TableCell>All Simples</TableCell>
            <TableCell>
              A Hand formed without Terminal or Honor Tiles.
            </TableCell>
            <TableCell>2</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>69</TableCell>
            <TableCell>Pure Double Chow</TableCell>
            <TableCell>Two identical chows in the same suit.</TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>70</TableCell>
            <TableCell>Mixed Double Chow</TableCell>
            <TableCell>
              Two chows of the same numbers but in different suits.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>71</TableCell>
            <TableCell>Short Straight</TableCell>
            <TableCell>
              Two chows in the same suit that run consecutively after one
              another to make a six-tile straight.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>72</TableCell>
            <TableCell>Two Terminal Chows</TableCell>
            <TableCell>Chows of 1-2-3 and 7-8-9 in the same suit.</TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>73</TableCell>
            <TableCell>Pung of Terminals or Honors</TableCell>
            <TableCell>
              A Pung or Kong of Ones, Nines, or Winds. (A dragon pung scores 2
              points.)
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>74</TableCell>
            <TableCell>Melded Kong</TableCell>
            <TableCell>
              A Kong that was claimed from another player or promoted from a
              melded pung.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>75</TableCell>
            <TableCell>One Voided Suit</TableCell>
            <TableCell>
              A Hand that uses tiles from only two of the three suits (it lacks
              any tiles from one of the three suits).
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>76</TableCell>
            <TableCell>No Honors</TableCell>
            <TableCell>
              A Hand formed entirely of suit tiles, without Winds or Dragons.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>77</TableCell>
            <TableCell>Edge Wait</TableCell>
            <TableCell>
              Waiting solely for a 3 to form a 1-2-3 chow, or solely for a 7 to
              form a 7-8-9 chow. Not valid if waiting for more than one tile.
              Not valid if the edge wait is combined with any other waits.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>78</TableCell>
            <TableCell>Closed Wait</TableCell>
            <TableCell>
              Waiting solely for a tile whose number is “inside” (in the middle)
              to form a chow. Not valid if waiting for more than one tile. Not
              valid if the closed wait is combined with other waits.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>79</TableCell>
            <TableCell>Single Wait</TableCell>
            <TableCell>
              Waiting solely for a tile to form a pair. Not valid if waiting for
              more than one tile (for example, holding 1-2-3-4 and waiting on
              the 1 and 4).
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>80</TableCell>
            <TableCell>Self-Drawn</TableCell>
            <TableCell>
              Going out (making mahjong) with a fresh tile picked from the wall.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>81</TableCell>
            <TableCell>Flower Tiles</TableCell>
            <TableCell>
              Each tile carved with Chinese word of Spring (or Summer, Autumn,
              Winder, Plum, Orchid, Bamboo, Chrysanthemum) will award you one
              point when you succeed in Hu. Flower replacement to be Hu, the
              point of Self-Drawn can add. It can’t add the points of Hu While
              Kong, it not Flower replacement can discard.
            </TableCell>
            <TableCell>1</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <div className="text-end my-3">
        <Link
          color="foreground"
          size="sm"
          isExternal
          showAnchorIcon
          isBlock
          href="fans.jpeg"
        >
          Same, but cooler.jpg
        </Link>
      </div>
    </div>
  );
}
