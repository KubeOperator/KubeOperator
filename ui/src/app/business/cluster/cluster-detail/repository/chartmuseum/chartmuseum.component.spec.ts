import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ChartmuseumComponent } from './chartmuseum.component';

describe('ChartmuseumComponent', () => {
  let component: ChartmuseumComponent;
  let fixture: ComponentFixture<ChartmuseumComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ChartmuseumComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ChartmuseumComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
