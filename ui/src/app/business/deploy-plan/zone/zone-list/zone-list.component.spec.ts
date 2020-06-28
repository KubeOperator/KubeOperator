import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ZoneListComponent } from './zone-list.component';

describe('ZoneListComponent', () => {
  let component: ZoneListComponent;
  let fixture: ComponentFixture<ZoneListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ZoneListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ZoneListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
